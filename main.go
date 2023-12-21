package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/joeychilson/hackernews/client"
	"github.com/joeychilson/hackernews/handlers"
	"github.com/joeychilson/hackernews/pages"
)

//go:embed assets/*
var assets embed.FS

func main() {
	client := client.New()

	mux := http.NewServeMux()

	// Redirects to GitHub repo
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			pages.NotFound().Render(r.Context(), w)
			return
		}
		http.Redirect(w, r, "/news", http.StatusFound)
	})

	// Assets
	mux.Handle("/assets/", http.FileServer(http.FS(assets)))

	// Pages
	mux.HandleFunc("/ask", handlers.HandleAsk(client))
	mux.HandleFunc("/item", handlers.HandleItem(client))
	mux.HandleFunc("/jobs", handlers.HandleJobs(client))
	mux.HandleFunc("/newest", handlers.HandleNewest(client))
	mux.HandleFunc("/news", handlers.HandleNews(client))
	mux.HandleFunc("/show", handlers.HandleShow(client))
	mux.HandleFunc("/submitted", handlers.HandleSubmitted(client))
	mux.HandleFunc("/threads", handlers.HandleThreads(client))
	mux.HandleFunc("/user", handlers.HandleUser(client))

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
