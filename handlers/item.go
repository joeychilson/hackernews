package handlers

import (
	"net/http"
	"strconv"

	"github.com/joeychilson/hackernews/pages"
	"github.com/joeychilson/hackernews/pkg/hackernews"
)

func HandleItem(c *hackernews.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			pages.Error().Render(r.Context(), w)
			return
		}

		item, err := c.GetItem(r.Context(), idInt)
		if err != nil {
			pages.NotFound().Render(r.Context(), w)
			return
		}

		if item.ID == 0 {
			pages.NotFound().Render(r.Context(), w)
			return
		}

		comments, err := c.Comments(r.Context(), item.Kids)
		if err != nil {
			pages.Error().Render(r.Context(), w)
			return
		}

		props := pages.ItemProps{
			Item:     item,
			Comments: comments,
		}

		pages.Item(props).Render(r.Context(), w)
	}
}
