package main

import (
	"context"
	rr "github.com/go-redis/redis_rate/v10"
	"github.com/google/uuid"
	"log"
	"net/http"
	"url-shortener/internal/database"
	"url-shortener/internal/response"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New()
		client, err := database.CreateClient()
		if err != nil {
			log.Printf("middleware: %s requestID, cannot reach redis server: %v\n", requestID, err)
			response.StatusInternalServerError(w)
			return
		}
		defer func() {
			err = client.Close()
			if err != nil {
				log.Printf("middleware: %s requestID: failed to close redis client: %v\n", requestID, err)
			}
		}()
		limiter := rr.NewLimiter(client)
		res, err := limiter.Allow(context.Background(), r.RemoteAddr, rr.PerSecond(10))
		if err != nil {
			log.Printf("middleware: %s requestID: redis limiter error: %v\n", requestID, err)
			response.StatusInternalServerError(w)
			return
		}
		if res.Remaining == 0 {
			log.Printf("middleware: %s requestID, %s, remote addr, too many requests\n", requestID, r.RemoteAddr)
			response.StatusTooManyRequestsResponse(w)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "requestID", requestID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
