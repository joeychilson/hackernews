package handlers

import (
	"net/http"
	"strconv"

	"github.com/joeychilson/hackernews/client"
	"github.com/joeychilson/hackernews/pages"
	"github.com/joeychilson/hackernews/types"
)

func HandleSubmitted(client *client.Client) http.HandlerFunc {
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

		var submitted []types.Item
		for _, id := range user.Submitted {
			item, err := client.Item(r.Context(), id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
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

		props := types.FeedProps{
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