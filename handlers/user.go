package handlers

import (
	"net/http"

	"github.com/joeychilson/hackernews/client"
	"github.com/joeychilson/hackernews/pages"
)

func HandleUser(c *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			pages.NotFound().Render(r.Context(), w)
			return
		}

		user, err := c.GetUser(r.Context(), id)
		if err != nil {
			pages.Error().Render(r.Context(), w)
			return
		}

		if user.ID == "" {
			pages.NotFound().Render(r.Context(), w)
			return
		}

		props := pages.UserProps{
			User: user,
		}

		pages.User(props).Render(r.Context(), w)
	}
}
