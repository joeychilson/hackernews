package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/joeychilson/hackernews/client"
	"github.com/joeychilson/hackernews/pages"
)

func HandleItem(c *client.Client) http.HandlerFunc {
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

		comments, err := getComments(r.Context(), c, item.Kids)
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

func getComments(ctx context.Context, c *client.Client, kids []int) ([]client.Item, error) {
	var comments []client.Item
	for _, kid := range kids {
		comment, err := c.GetItem(ctx, kid)
		if err != nil {
			continue
		}
		if len(comment.Kids) > 0 {
			children, err := getComments(ctx, c, comment.Kids)
			if err != nil {
				continue
			}
			comment.Children = children
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
