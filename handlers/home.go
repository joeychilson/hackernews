package handlers

import (
	"net/http"
	"strconv"

	"github.com/joeychilson/hackernews/client"
	"github.com/joeychilson/hackernews/models"
	"github.com/joeychilson/hackernews/pages"
)

const (
	pageSize     = 30
	visiblePages = 5
)

func HomePage(client *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("p")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		topStoryIDs, err := client.TopStories(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		start := (page - 1) * pageSize
		end := start + pageSize
		if start > len(topStoryIDs) {
			start = len(topStoryIDs)
		}
		if end > len(topStoryIDs) {
			end = len(topStoryIDs)
		}

		topStories := make([]models.Item, 0, pageSize)
		for _, id := range topStoryIDs[start:end] {
			story, err := client.GetItem(r.Context(), id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			topStories = append(topStories, story)
		}

		totalPages := len(topStoryIDs)/pageSize + 1

		startPage := max(1, page-(visiblePages/2))
		if startPage+visiblePages > totalPages {
			startPage = max(1, totalPages-visiblePages+1)
		}

		endPage := min(startPage+visiblePages-1, totalPages)

		pageNumbers := make([]int, 0, endPage-startPage+1)
		for i := startPage; i <= endPage; i++ {
			pageNumbers = append(pageNumbers, i)
		}

		props := models.HomeProps{
			Stories:     topStories,
			Total:       len(topStoryIDs),
			PerPage:     pageSize,
			CurrentPage: page,
			StartPage:   startPage,
			TotalPages:  totalPages,
			PageNumbers: pageNumbers,
		}
		pages.Home(props).Render(r.Context(), w)
	}
}
