package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Title       string    `xml:"channel>title"`
	Description string    `xml:"channel>description"`
	Link        string    `xml:"channel>link"`
	Items       []RSSItem `xml:"channel>item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
}

// fetchFeed fetches an RSS feed from the given URL and parses it into an RSSFeed
// struct. It also unescapes HTML entities in the feed fields.
//
// The function uses the given context to cancel the HTTP request if it
// times out or is canceled.
//
// If the HTTP request fails, the function returns an error. If the HTTP
// request succeeds but the response body is not valid XML, the function
// returns an error.
//
// If the function succeeds, it returns a pointer to the parsed RSSFeed
// struct.
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Create an HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent header
	req.Header.Add("User-Agent", "gator")

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected HTTP status: " + resp.Status)
	}

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read feed data: %w", err)
	}

	// Parse the XML into the RSSFeed struct
	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	// Decode escaped HTML entities in the feed fields
	feed.Title = html.UnescapeString(feed.Title)
	feed.Description = html.UnescapeString(feed.Description)
	for i, item := range feed.Items {
		feed.Items[i].Title = html.UnescapeString(item.Title)
		feed.Items[i].Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}
