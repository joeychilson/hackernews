package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joeychilson/hackernews/client"
	"github.com/joeychilson/hackernews/handlers"
)

//go:embed dist/*
var dist embed.FS

func main() {
	client := client.New()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Handle("/dist/*", http.FileServer(http.FS(dist)))
	router.HandleFunc("/", handlers.HomePage(client))

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
