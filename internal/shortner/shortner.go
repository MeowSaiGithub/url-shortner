package shortner

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/cache/v9"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"url-shortener/internal/database"
	"url-shortener/internal/model"
	"url-shortener/internal/response"
)

const (
	alphabet      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	defaultExpiry = 1 * time.Hour
)

func base62Encode(number uint64) string {
	length := len(alphabet)
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(10)
	for ; number > 0; number = number / uint64(length) {
		encodedBuilder.WriteByte(alphabet[(number % uint64(length))])
	}
	return encodedBuilder.String()
}

func Shorten(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("requestID")
	message := model.Request{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		log.Printf("shorten: %s requestID, failed to decode request: %v\n", requestID, err)
		response.StatusBadRequestResponse(w, "bad request")
		return
	}

	if !strings.HasPrefix(message.URL, "http://") && !strings.HasPrefix(message.URL, "https://") {
		log.Printf("shorten: %s requestID, URL must be with either http:// or https://\n", requestID)
		response.StatusBadRequestResponse(w, "URL must be with either http:// or https://")
		return
	}

	var id string
	if message.CustomShort == "" {
		id = base62Encode(rand.Uint64())
	} else {
		id = message.CustomShort
	}
	var expiry time.Duration
	if message.Expiry == 0 {
		expiry = defaultExpiry
	} else {
		expiry = message.Expiry
	}
	client, err := database.CreateClient()
	if err != nil {
		log.Printf("shorten: %s requestID, cannot reach redis server: %v\n", requestID, err)
		response.StatusInternalServerError(w)
		return
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("shorten: %s requestID, failed to close client: %v\n", requestID, err)
		}
	}()
	rdb := cache.New(&cache.Options{
		Redis: client,
	})
	ctx := context.Background()
	exist := rdb.Exists(ctx, id)
	if exist {
		log.Printf("shorten: %s requestID, %s URL already exists\n", requestID, id)
		response.StatusConflictedResponse(w, "URL already in-used, try again later.")
		return
	}
	item := &cache.Item{
		Ctx:   ctx,
		Key:   id,
		Value: message.URL,
		TTL:   expiry,
	}
	if err := rdb.Set(item); err != nil {
		log.Printf("shorten: %s requestID, failed to set cache: %v\n", err)
		response.StatusInternalServerError(w)
		return
	}
	response.StatusOKResponse(w, id, expiry)
}

func Resolve(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("requestID")
	vars := mux.Vars(r)
	url := vars["id"]
	if url == "" {
		log.Printf("resolve: %s requestID, URL is empty", requestID)
		response.StatusBadRequestResponse(w, "URL is empty")
		return
	}
	client, err := database.CreateClient()
	if err != nil {
		log.Printf("resolve: %s requestID, cannot reach redis server: %v\n", requestID, err)
		response.StatusInternalServerError(w)
		return
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("resolve:  %s requestID, failed to close client: %v\n", requestID, err)
		}
	}()
	rdb := cache.New(&cache.Options{
		Redis: client,
	})
	ctx := context.Background()
	var redirect string
	if err := rdb.Get(ctx, url, &redirect); err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			log.Printf("resolve: %s requestID, %s key does not exist\n", requestID, url)
			response.StatusNotFoundResponse(w, "")
			return
		}
		log.Printf("resolve: %s requestID, failed to key cache: %v\n", requestID, err)
		response.StatusInternalServerError(w)
		return
	}
	response.StatusRedirectResponse(w, r, redirect)
}
