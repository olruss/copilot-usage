package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	baseURL    = "https://api.github.com"
	apiVersion = "2022-11-28"
	perPage    = 100
)

// Client is a GitHub API client for Copilot metrics.
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient creates a new API client with the given auth token.
func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{},
	}
}

// FetchMetrics fetches Copilot usage metrics for the given org and date range.
// since and until are optional ISO 8601 date strings (YYYY-MM-DD).
func (c *Client) FetchMetrics(org, since, until string) ([]DayMetrics, error) {
	var allMetrics []DayMetrics
	page := 1

	for {
		u, err := url.Parse(fmt.Sprintf("%s/orgs/%s/copilot/metrics", baseURL, org))
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}

		q := u.Query()
		q.Set("per_page", fmt.Sprintf("%d", perPage))
		q.Set("page", fmt.Sprintf("%d", page))
		if since != "" {
			q.Set("since", since)
		}
		if until != "" {
			q.Set("until", until)
		}
		u.RawQuery = q.Encode()

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}

		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("X-GitHub-Api-Version", apiVersion)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("API request failed: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("reading response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
		}

		var metrics []DayMetrics
		if err := json.Unmarshal(body, &metrics); err != nil {
			return nil, fmt.Errorf("parsing response: %w", err)
		}

		allMetrics = append(allMetrics, metrics...)

		// If we got fewer results than per_page, we've reached the last page
		if len(metrics) < perPage {
			break
		}
		page++
	}

	return allMetrics, nil
}
