package handlers

import (
	"net/http"
	"strconv"

	"github.com/joeychilson/hackernews/client"
	"github.com/joeychilson/hackernews/pages"
)

func HandleJobs(c *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("p")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		storyIDs, err := c.GetJobsStories(r.Context())
		if err != nil {
			pages.Error().Render(r.Context(), w)
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

		stories := make([]client.Item, 0, pageSize)
		for _, id := range storyIDs[start:end] {
			story, err := c.GetItem(r.Context(), id)
			if err != nil {
				pages.NotFound().Render(r.Context(), w)
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

		props := pages.FeedProps{
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
