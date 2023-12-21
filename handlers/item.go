package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/joeychilson/hackernews/client"
	"github.com/joeychilson/hackernews/pages"
	"github.com/joeychilson/hackernews/types"
)

func HandleItem(client *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			pages.NotFound().Render(r.Context(), w)
			return
		}

		item, err := client.GetItem(r.Context(), idInt)
		if err != nil {
			pages.NotFound().Render(r.Context(), w)
			return
		}

		comments, err := getComments(r.Context(), client, item.Kids)
		if err != nil {
			pages.NotFound().Render(r.Context(), w)
			return
		}

		props := types.ItemProps{
			Item:     item,
			Comments: comments,
		}

		pages.Item(props).Render(r.Context(), w)
	}
}

func getComments(ctx context.Context, client *client.Client, kids []int) ([]types.Item, error) {
	var comments []types.Item
	for _, kid := range kids {
		comment, err := client.GetItem(ctx, kid)
		if err != nil {
			continue
		}
		if len(comment.Kids) > 0 {
			children, err := getComments(ctx, client, comment.Kids)
			if err != nil {
				continue
			}
			comment.Children = children
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
