package handlers

import (
	"net/http"

	"github.com/joeychilson/hackernews/client"
	"github.com/joeychilson/hackernews/pages"
	"github.com/joeychilson/hackernews/types"
)

func HandleUser(client *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		user, err := client.User(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		props := types.UserProps{
			User: user,
		}

		pages.User(props).Render(r.Context(), w)
	}
}
