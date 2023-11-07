package main

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
	"url-shortener/internal/shortner"
)

func (a *app) setRouter(path string, port string) {
	a.router = mux.NewRouter()
	a.router.Use(middleware)

	a.router.HandleFunc(path+"generate", shortner.Shorten).Methods(http.MethodPost)
	a.router.HandleFunc(path+"{id}", shortner.Resolve).Methods(http.MethodGet)
	a.srv = &http.Server{
		Addr:         ":" + port,
		Handler:      a.router,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
}

func (a *app) run() error {
	return a.srv.ListenAndServe()
}

func (a *app) shutdown(ctx context.Context, cancel context.CancelFunc, sig chan os.Signal) {
	<-sig
	defer func() {
		cancel()
		log.Println("server exited properly.")
	}()
	go func() {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Fatal("graceful shutdown timed out.. forcing exit.")
		}
	}()
	if err := a.srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %+v", err)
	}
}
