package handlers

import (
	"net/http"
	"strconv"

	"github.com/joeychilson/hackernews/client"
	"github.com/joeychilson/hackernews/pages"
	"github.com/joeychilson/hackernews/types"
)

func HandleJobs(client *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("p")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		storyIDs, err := client.JobsStories(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		start := (page - 1) * pageSize
		end := start + pageSize
		if start > len(storyIDs) {
			start = len(storyIDs)
		}
		if end > len(storyIDs) {
			end = len(storyIDs)
		}

		stories := make([]types.Item, 0, pageSize)
		for _, id := range storyIDs[start:end] {
			story, err := client.Item(r.Context(), id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			stories = append(stories, story)
		}

		totalPages := len(stories)/pageSize + 1

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
			Title:       "Jobs",
			Stories:     stories,
			Total:       len(storyIDs),
			PerPage:     pageSize,
			CurrentPage: page,
			StartPage:   startPage,
			TotalPages:  totalPages,
			PageNumbers: pageNumbers,
		}
		pages.Feed(props).Render(r.Context(), w)
	}
}
