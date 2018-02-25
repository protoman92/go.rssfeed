package gorssfeed

import (
	"errors"
	"testing"
	"time"

	"github.com/SlyMarbo/rss"
)

const (
	sleepTime = 1e9
)

type mockFetcher struct {
	mockFetch func() *FetchResult
}

func (mf *mockFetcher) Fetch() *FetchResult {
	return mf.mockFetch()
}

func newErrorFetcher(err error) *mockFetcher {
	return &mockFetcher{mockFetch: func() *FetchResult {
		return &FetchResult{nil, err}
	}}
}

func newItemFetcher(items []*rss.Item) *mockFetcher {
	return &mockFetcher{mockFetch: func() *FetchResult {
		return &FetchResult{feed: &rss.Feed{Items: items}, err: nil}
	}}
}

func Test_SubscriptionWithFetchError_ShouldRelayOnClose(t *testing.T) {
	// Setup
	t.Parallel()
	err := errors.New("Fetch error")
	fetcher := newErrorFetcher(err)

	// When
	sub := Subscribe(fetcher)
	time.Sleep(sleepTime)

	// Then
	err1 := sub.Unsubscribe()

	if err1 == nil {
		t.Errorf("Should have %v, but instead got nothing", err)
	} else if err1.Error() != err.Error() {
		t.Errorf("Should have %v, but instead got %v", err, err1)
	}
}

func Test_FetchWithNoItems_ShouldNotUpdate(t *testing.T) {
	// Setup
	t.Parallel()
	fetcher := newItemFetcher([]*rss.Item{})

	// When
	sub := Subscribe(fetcher)
	time.Sleep(sleepTime)

	// Then
	length := len(sub.Updates())

	if length != 0 {
		t.Errorf("Should not have any update, but instead got %d", length)
	}
}
