package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/joeychilson/hackernews/handlers"
	"github.com/joeychilson/hackernews/pages"
	"github.com/joeychilson/hackernews/pkg/hackernews"
)

//go:embed static
var static embed.FS

func main() {
	client := hackernews.New()

	mux := http.NewServeMux()

	// Redirects to News to match the original site
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			pages.NotFound().Render(r.Context(), w)
			return
		}
		http.Redirect(w, r, "/news", http.StatusFound)
	})

	// Static files
	mux.Handle("/static/", http.FileServer(http.FS(static)))

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

	log.Println("Listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
