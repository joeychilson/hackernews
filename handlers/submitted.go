package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

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

		user, err := c.User(r.Context(), id)
		if err != nil {
			pages.Error().Render(r.Context(), w)
			return
		}

		if user.ID == "" {
			pages.NotFound().Render(r.Context(), w)
			return
		}

		var (
			wg           sync.WaitGroup
			errOnce      sync.Once
			submittedMap = make(map[int]*hackernews.Item)
			submitted    []*hackernews.Item
		)
		mu := &sync.Mutex{}

		for _, id := range user.Submitted {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				item, e := c.Item(r.Context(), id)
				if e != nil {
					errOnce.Do(func() {
						err = e
					})
					return
				}
				if item.Type == "story" {
					mu.Lock()
					submittedMap[id] = item
					mu.Unlock()
				}
			}(id)
		}
		wg.Wait()

		// We need to sort the submitted stories by the order they were submitted
		for _, id := range user.Submitted {
			if item, exists := submittedMap[id]; exists {
				submitted = append(submitted, item)
			}
		}

		pageStr := r.URL.Query().Get("p")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		start := (page - 1) * pageSize
		end := start + pageSize
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
