package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/joeychilson/hackernews/pages"
)

//go:embed dist/*
var dist embed.FS

func main() {
	mux := http.NewServeMux()

	mux.Handle("/dist/", http.FileServer(http.FS(dist)))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pages.Home().Render(r.Context(), w)
	})

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
