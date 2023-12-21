package hackernews

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultURL = "https://hacker-news.firebaseio.com/v0/"

type Client struct {
	baseURL string
	client  *http.Client
}

func New() *Client {
	return &Client{
		baseURL: defaultURL,
		client:  http.DefaultClient,
	}
}

type Item struct {
	ID          int    `json:"id,omitempty"`
	Deleted     bool   `json:"deleted,omitempty"`
	Type        string `json:"type,omitempty"`
	By          string `json:"by,omitempty"`
	Time        int64  `json:"time,omitempty"`
	Text        string `json:"text,omitempty"`
	Dead        bool   `json:"dead,omitempty"`
	Parent      int    `json:"parent,omitempty"`
	Kids        []int  `json:"kids,omitempty"`
	URL         string `json:"url,omitempty"`
	Score       int    `json:"score,omitempty"`
	Title       string `json:"title,omitempty"`
	Parts       []int  `json:"parts,omitempty"`
	Descendants int    `json:"descendants,omitempty"`
	Children    []Item `json:"children,omitempty"`
}

type User struct {
	ID        string `json:"id"`
	Created   int    `json:"created"`
	Karma     int    `json:"karma"`
	About     string `json:"about"`
	Submitted []int  `json:"submitted"`
}

func (c *Client) GetAskStories(ctx context.Context) ([]int, error) {
	b, err := c.get(ctx, "askstories.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get ask stories: %w", err)
	}
	var stories []int
	err = json.Unmarshal(b, &stories)
	return stories, err
}

func (c *Client) GetJobsStories(ctx context.Context) ([]int, error) {
	b, err := c.get(ctx, "jobstories.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get job stories: %w", err)
	}
	var stories []int
	err = json.Unmarshal(b, &stories)
	return stories, err
}

func (c *Client) GetNewestStories(ctx context.Context) ([]int, error) {
	b, err := c.get(ctx, "newstories.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get newest stories: %w", err)
	}
	var stories []int
	err = json.Unmarshal(b, &stories)
	return stories, err
}

func (c *Client) GetShowStories(ctx context.Context) ([]int, error) {
	b, err := c.get(ctx, "showstories.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get show stories: %w", err)
	}
	var stories []int
	err = json.Unmarshal(b, &stories)
	return stories, err
}

func (c *Client) GetTopStories(ctx context.Context) ([]int, error) {
	b, err := c.get(ctx, "topstories.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get top stories: %w", err)
	}
	var stories []int
	err = json.Unmarshal(b, &stories)
	return stories, err
}

func (c *Client) GetItem(ctx context.Context, id int) (Item, error) {
	path := fmt.Sprintf("item/%d.json", id)
	b, err := c.get(ctx, path)
	if err != nil {
		return Item{}, fmt.Errorf("failed to get item: %w", err)
	}
	item := Item{}
	err = json.Unmarshal(b, &item)
	return item, err
}

func (c *Client) GetUser(ctx context.Context, id string) (User, error) {
	path := fmt.Sprintf("/user/%s.json", id)
	b, err := c.get(ctx, path)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	user := User{}
	err = json.Unmarshal(b, &user)
	return user, err
}

func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
