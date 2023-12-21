package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joeychilson/hackernews/client"
	"github.com/joeychilson/hackernews/handlers"
	"github.com/joeychilson/hackernews/pages"
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

	// Assets
	router.Handle("/dist/*", http.FileServer(http.FS(dist)))

	// Pages
	router.HandleFunc("/", handlers.HomePage(client))
	router.HandleFunc("/item", handlers.ItemPage(client))

	// Not Found
	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		pages.NotFound().Render(r.Context(), w)
	})

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
