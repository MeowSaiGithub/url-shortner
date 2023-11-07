package main

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultPath = "/"
	defaultPort = "8080"
)

type app struct {
	router *mux.Router
	srv    *http.Server
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load environment file.")
	}

	path := os.Getenv("APP_PATH")
	if path == "" {
		path = defaultPath
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = defaultPort
	}

	a := app{}
	a.setRouter(path, port)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	log.Printf("basic path: %s\n", path)
	log.Printf("listening at port :%s\n", port)

	go func() {
		if err := a.run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v\n", err)
		}
	}()
	a.shutdown(ctx, cancel, sig)
}
