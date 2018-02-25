package gorssfeed

import (
	"github.com/SlyMarbo/rss"
)

// FetchResult represents the result of a fetch.
type FetchResult struct {
	feed *rss.Feed
	err  error
}

// Fetcher represents an object that can fetch remote data.
type Fetcher interface {
	Fetch() *FetchResult
}

type fetcher struct {
	url string
}

func (f *fetcher) Fetch() *FetchResult {
	feed, err := rss.Fetch(f.url)
	return &FetchResult{feed: feed, err: err}
}

// NewFetcher creates a new fetcher.
func NewFetcher(url string) Fetcher {
	return &fetcher{url: url}
}
