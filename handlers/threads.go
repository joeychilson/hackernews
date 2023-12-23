package handlers

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/joeychilson/hackernews/pages"
	"github.com/joeychilson/hackernews/pkg/hackernews"
)

func HandleThreads(c *hackernews.Client) http.HandlerFunc {
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
			wg         sync.WaitGroup
			errOnce    sync.Once
			threadsMap = make(map[int]*hackernews.Item)
			threads    []*hackernews.Item
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
				if item.Type == "comment" {
					mu.Lock()
					threadsMap[id] = item
					mu.Unlock()
				}
			}(id)
		}
		wg.Wait()

		// We need to sort the submitted stories by the order they were submitted
		for _, id := range user.Submitted {
			if item, exists := threadsMap[id]; exists {
				threads = append(threads, item)
			}
		}

		pageStr := r.URL.Query().Get("p")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		start := (page - 1) * pageSize
		end := start + pageSize
		if start > len(threads) {
			start = len(threads)
		}
		if end > len(threads) {
			end = len(threads)
		}
		threads = threads[start:end]

		threadItems := make([]*hackernews.Item, 0, len(threads))
		for _, thread := range threads {
			item, err := c.Item(r.Context(), thread.ID)
			if err != nil {
				pages.NotFound().Render(r.Context(), w)
				return
			}

			comments, err := c.Comments(r.Context(), item.Kids)
			if err != nil {
				pages.Error().Render(r.Context(), w)
				return
			}
			item.Children = comments

			threadItems = append(threadItems, item)
		}

		totalPages := len(threads)/pageSize + 1

		startPage := max(1, page-(visiblePages/2))
		if startPage+visiblePages > totalPages {
			startPage = max(1, totalPages-visiblePages+1)
		}

		endPage := min(startPage+visiblePages-1, totalPages)

		pageNumbers := make([]int, 0, endPage-startPage+1)
		for i := startPage; i <= endPage; i++ {
			pageNumbers = append(pageNumbers, i)
		}

		props := pages.UserCommentsProps{
			User:        user.ID,
			Comments:    threadItems,
			Total:       len(threadItems),
			PerPage:     pageSize,
			CurrentPage: page,
			StartPage:   startPage,
			TotalPages:  totalPages,
			PageNumbers: pageNumbers,
		}
		pages.UserComments(props).Render(r.Context(), w)
	}
}
