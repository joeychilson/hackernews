package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/joeychilson/hackernews/pages"
	"github.com/joeychilson/hackernews/pkg/hackernews"
)

func HandleSubmitted(c *hackernews.Client) http.HandlerFunc {
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

		var submitted []hackernews.Item
		for _, id := range user.Submitted {
			item, err := c.GetItem(r.Context(), id)
			if err != nil {
				pages.NotFound().Render(r.Context(), w)
				return
			}
			if item.Type != "story" {
				continue
			}
			submitted = append(submitted, item)
		}

		pageStr := r.URL.Query().Get("p")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		start := (page - 1) * pageSize
		end := start + pageSize
		if start > len(submitted) {
			start = len(submitted)
		}
		if end > len(submitted) {
			end = len(submitted)
		}

		submitted = submitted[start:end]

		totalPages := len(submitted)/pageSize + 1

		startPage := max(1, page-(visiblePages/2))
		if startPage+visiblePages > totalPages {
			startPage = max(1, totalPages-visiblePages+1)
		}

		endPage := min(startPage+visiblePages-1, totalPages)

		pageNumbers := make([]int, 0, endPage-startPage+1)
		for i := startPage; i <= endPage; i++ {
			pageNumbers = append(pageNumbers, i)
		}

		props := pages.FeedProps{
			Title:       fmt.Sprintf("%s's submissions", user.ID),
			Stories:     submitted,
			Total:       len(submitted),
			PerPage:     pageSize,
			CurrentPage: page,
			StartPage:   startPage,
			TotalPages:  totalPages,
			PageNumbers: pageNumbers,
		}
		pages.Feed(props).Render(r.Context(), w)
	}
}
