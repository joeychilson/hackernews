package handlers

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/joeychilson/hackernews/pages"
	"github.com/joeychilson/hackernews/pkg/hackernews"
)

func HandleNewest(c *hackernews.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("p")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		storyIDs, err := c.GetNewestStories(r.Context())
		if err != nil {
			pages.Error().Render(r.Context(), w)
			return
		}

		start := (page - 1) * pageSize
		end := start + pageSize
		if end > len(storyIDs) {
			end = len(storyIDs)
		}

		paginatedIDs := storyIDs[start:end]

		var wg sync.WaitGroup

		stories := make([]hackernews.Item, len(paginatedIDs))
		for i, id := range paginatedIDs {
			wg.Add(1)
			go func(i, id int) {
				defer wg.Done()

				story, err := c.GetItem(r.Context(), id)
				if err != nil {
					return
				}

				stories[i] = story
			}(i, id)
		}
		wg.Wait()

		totalPages := len(storyIDs)/pageSize + 1

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
