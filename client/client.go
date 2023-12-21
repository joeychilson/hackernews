package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/joeychilson/hackernews/types"
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

func (c *Client) AskStories(ctx context.Context) ([]int, error) {
	b, err := c.get(ctx, "askstories.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get ask stories: %w", err)
	}
	var stories []int
	err = json.Unmarshal(b, &stories)
	return stories, err
}

func (c *Client) JobsStories(ctx context.Context) ([]int, error) {
	b, err := c.get(ctx, "jobstories.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get job stories: %w", err)
	}
	var stories []int
	err = json.Unmarshal(b, &stories)
	return stories, err
}

func (c *Client) NewestStories(ctx context.Context) ([]int, error) {
	b, err := c.get(ctx, "newstories.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get newest stories: %w", err)
	}
	var stories []int
	err = json.Unmarshal(b, &stories)
	return stories, err
}

func (c *Client) ShowStories(ctx context.Context) ([]int, error) {
	b, err := c.get(ctx, "showstories.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get show stories: %w", err)
	}
	var stories []int
	err = json.Unmarshal(b, &stories)
	return stories, err
}

func (c *Client) TopStories(ctx context.Context) ([]int, error) {
	b, err := c.get(ctx, "topstories.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get top stories: %w", err)
	}
	var stories []int
	err = json.Unmarshal(b, &stories)
	return stories, err
}

func (c *Client) GetItem(ctx context.Context, id int) (types.Item, error) {
	path := fmt.Sprintf("item/%d.json", id)
	b, err := c.get(ctx, path)
	if err != nil {
		return types.Item{}, fmt.Errorf("while getting item: %w", err)
	}
	item := types.Item{}
	err = json.Unmarshal(b, &item)
	return item, err
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
