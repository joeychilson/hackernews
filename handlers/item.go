package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"

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

func getComments(ctx context.Context, c *hackernews.Client, kids []int) ([]hackernews.Item, error) {
	var wg sync.WaitGroup
	commentCh := make(chan hackernews.Item, len(kids))
	errCh := make(chan error, len(kids))

	for _, kid := range kids {
		wg.Add(1)
		go func(kid int) {
			defer wg.Done()
			comment, err := c.GetItem(ctx, kid)
			if err != nil {
				errCh <- err
				return
			}
			if len(comment.Kids) > 0 {
				children, err := getComments(ctx, c, comment.Kids)
				if err != nil {
					errCh <- err
					return
				}
				comment.Children = children
			}
			commentCh <- comment
		}(kid)
	}

	wg.Wait()
	close(commentCh)
	close(errCh)

	if len(errCh) > 0 {
		return nil, fmt.Errorf("errors occurred while fetching comments")
	}

	var comments []hackernews.Item
	for comment := range commentCh {
		comments = append(comments, comment)
	}
	return comments, nil
}
