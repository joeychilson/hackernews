package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/joeychilson/hackernews/client"
	"github.com/joeychilson/hackernews/models"
	"github.com/joeychilson/hackernews/pages"
)

//go:embed dist/*
var dist embed.FS

func main() {
	client := client.New()

	mux := http.NewServeMux()

	mux.Handle("/dist/", http.FileServer(http.FS(dist)))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		topStoryIDs, err := client.TopStories(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		topStories := make([]models.Item, 0, 30)

		for _, id := range topStoryIDs[:30] {
			story, err := client.GetItem(r.Context(), id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			topStories = append(topStories, story)
		}

		pages.Home(topStories).Render(r.Context(), w)
	})

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
