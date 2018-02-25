package gorssfeed

import (
	"time"

	"github.com/SlyMarbo/rss"
)

// Subscription represents a RSS feed subscription with updates.
type Subscription interface {
	Updates() chan *rss.Item
	Unsubscribe() error
}

type subscription struct {
	fetcher Fetcher

	// This is a channel of channels. When we unsubscribe from this subscription,
	// send a closing request (in the form of a channel), and access whatever
	// error that is sent through this channel to return it.
	closing chan chan error

	updates chan *rss.Item
}

func (s *subscription) Updates() chan *rss.Item {
	return s.updates
}

func (s *subscription) Unsubscribe() error {
	errc := make(chan error)
	s.closing <- errc
	return <-errc
}

func (s *subscription) loop() {
	items := make([]*rss.Item, 0)
	var currentTime, refreshTime time.Time
	var fetchDelay time.Duration
	var err error
	var errc chan error
	var fetchDone chan *FetchResult

	for {
		var first *rss.Item
		var updates chan *rss.Item

		if len(items) > 0 {
			first = items[0]
			updates = s.updates
		}

		currentTime = time.Now()

		if refreshTime.After(currentTime) {
			fetchDelay = refreshTime.Sub(currentTime)
		} else {
			fetchDelay = 0
		}

		startFetch := time.After(fetchDelay)

		select {
		case <-startFetch:
			fetchDone = make(chan *FetchResult)

			go func() {
				fetchDone <- s.fetcher.Fetch()
			}()

		case updates <- first:
			// If items is empty, updates is a nil channel, which always blocks. Since
			// select will never choose a blocking channel, we can be sure that this
			// code is not executed when there are no items to send. Therefore, we
			// do not need to worry about out of index exception.
			items = items[1:]

		case result := <-fetchDone:
			fetchDone = nil

			if err = result.err; err != nil {
				break
			} else {
				items = append(items, result.feed.Items...)
				refreshTime = result.feed.Refresh
			}

		case errc = <-s.closing:
			// Error channels only appear when we unsubscribe from the subscription.
			errc <- err
			close(s.updates)
			return
		}
	}
}

// Subscribe to a RSS feed to receive updates.
func Subscribe(fetcher Fetcher) Subscription {
	s := subscription{
		fetcher: fetcher,
		closing: make(chan chan error),
		updates: make(chan *rss.Item),
	}

	go s.loop()
	return &s
}
