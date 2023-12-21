package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/joeychilson/hackernews/client"
	"github.com/joeychilson/hackernews/handlers"
)

//go:embed dist/*
var dist embed.FS

func main() {
	client := client.New()

	mux := http.NewServeMux()

	mux.Handle("/dist/", http.FileServer(http.FS(dist)))
	mux.HandleFunc("/", handlers.HandleHome(client))

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
