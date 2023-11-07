package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type errResp struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"ts"`
}

type okResp struct {
	URL       string        `json:"url"`
	Expire    time.Duration `json:"expiry"`
	Timestamp time.Time     `json:"timestamp"`
}

func StatusRedirectResponse(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}

func StatusInternalServerError(w http.ResponseWriter) {
	payload := errResp{
		Message:   "internal server error",
		Timestamp: time.Now(),
	}
	respond, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(respond)
}

func StatusConflictedResponse(w http.ResponseWriter, payload string) {
	resp := errResp{
		Message:   payload,
		Timestamp: time.Now(),
	}
	respond, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	_, _ = w.Write(respond)
}

func StatusTooManyRequestsResponse(w http.ResponseWriter) {
	resp := errResp{
		Message:   "too many requests, try again after a few minute",
		Timestamp: time.Now(),
	}
	respond, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	_, _ = w.Write(respond)
}

func StatusBadRequestResponse(w http.ResponseWriter, payload string) {
	resp := errResp{
		Message:   payload,
		Timestamp: time.Now(),
	}
	respond, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(respond)
}

func StatusOKResponse(w http.ResponseWriter, payload string, expiry time.Duration) {
	resp := okResp{
		URL:       payload,
		Expire:    expiry,
		Timestamp: time.Now(),
	}
	respond, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respond)
}

func StatusNotFoundResponse(w http.ResponseWriter, payload string) {
	resp := errResp{
		Message:   fmt.Sprintf("%s URL not found", payload),
		Timestamp: time.Now(),
	}
	respond, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write(respond)
}
