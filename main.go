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

//go:embed assets/*
var assets embed.FS

func main() {
	client := client.New()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Redirects to GitHub repo
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://www.github.com/joeychilson/hackernews", http.StatusFound)
	})

	// Assets
	router.Handle("/assets/*", http.FileServer(http.FS(assets)))

	// Pages
	router.HandleFunc("/ask", handlers.HandleAsk(client))
	router.HandleFunc("/item", handlers.HandleItem(client))
	router.HandleFunc("/jobs", handlers.HandleJobs(client))
	router.HandleFunc("/newest", handlers.HandleNewest(client))
	router.HandleFunc("/news", handlers.HandleNews(client))
	router.HandleFunc("/show", handlers.HandleShow(client))
	router.HandleFunc("/submitted", handlers.HandleSubmitted(client))
	router.HandleFunc("/threads", handlers.HandleThreads(client))
	router.HandleFunc("/user", handlers.HandleUser(client))

	// Not Found
	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		pages.NotFound().Render(r.Context(), w)
	})

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
