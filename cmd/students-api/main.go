package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shankar7042/students-api/internal/config"
)

func main() {
	// load config

	cfg := config.MustLoad()

	// database setup
	// setup router

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to golang crud api"))
	})
	// setup server

	server := http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: router,
	}

	fmt.Printf("Server started at %s", cfg.HttpServer.Addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("failed to start server")
	}
}
