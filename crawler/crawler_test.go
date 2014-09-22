package crawler

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"testing"

	// log "github.com/cihub/seelog"
	"github.com/davegardnerisme/deephash"
	"github.com/stretchr/testify/assert"

	"github.com/mattheath/kraken/domain"
)

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body   string
	urls   []string
	assets []string
}

func (f fakeFetcher) Fetch(target *url.URL) ([]*url.URL, []*url.URL, error) {
	if res, ok := f[target.String()]; ok {
		furls, _ := stringsToUrls(res.urls)
		fassets, _ := stringsToUrls(res.assets)
		return furls, fassets, nil
	}
	return nil, nil, errors.New("not found: " + target.String())
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
		[]string{},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
		[]string{},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
		[]string{},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
		[]string{},
	},
}

func TestCrawlSuccess(t *testing.T) {

	c := newMockCrawler()

	// Test completable requests
	testCases := map[string]*fakeResult{
		"http://golang.org/":     fetcher["http://golang.org/"],
		"http://golang.org/pkg/": fetcher["http://golang.org/pkg/"],
	}

	// Fire!
	for target, res := range testCases {
		c.crawl(strToUrl(target), 1, fetcher)

		var r *Result
		select {
		case <-c.skipped:
			t.Error("Request was skipped")
		case <-c.errored:
			t.Error("Request errored")
		case r = <-c.completed:
		}

		assert.NotNil(t, r.Page)
		assert.NotNil(t, r.Page.Links)

		// To compare links returned we need to ensure they are in the same datastructure
		// and in the same order. Converting to a sorted string slice ensures this.
		links := linksToStrings(r.Page.Links)
		tclinks := res.urls
		sort.Strings(links)
		sort.Strings(tclinks)

		// Hash both to do a deep comparision of the two
		assert.Equal(t, deephash.Hash(links), deephash.Hash(tclinks), fmt.Sprintf("%#v and %#v should be equal", links, tclinks))
	}

}

func TestCrawlError(t *testing.T) {

	c := newMockCrawler()

	// Test error cases, these don't exist in the mocked fetcher
	testCases := []string{
		"http://www.omfgdogs.com/",
		"http://ducksarethebest.com/",
	}

	// Fire!
	for _, target := range testCases {
		c.crawl(strToUrl(target), 1, fetcher)

		var r *Result
		select {
		case <-c.completed:
			t.Error("Request completed, should have errored")
		case <-c.skipped:
			t.Error("Request skipped, should have errored")
		case r = <-c.errored:
		}

		assert.NotNil(t, r.Error)
	}

}

func TestCrawlSkipAtMaxDepth(t *testing.T) {

	c := newMockCrawler()

	// Test error cases, these don't exist in the mocked fetcher
	testCases := []string{
		"http://golang.org/pkg/",
	}

	// Fire!
	for _, target := range testCases {
		c.crawl(strToUrl(target), 0, fetcher)

		var r *Result
		select {
		case <-c.completed:
			t.Error("Request completed, should have been skipped")
		case <-c.errored:
			t.Error("Request errored, should have been skipped")
		case r = <-c.skipped:
		}

		assert.NotNil(t, r)
	}

}

// newMockCrawler returns a crawler with buffered channels
// suitable for single threaded use
func newMockCrawler() *crawler {
	c := NewCrawler()

	// Mock channels as buffered so we can retrieve our results in a single thread
	c.skipped = make(chan *Result, 1)
	c.completed = make(chan *Result, 1)
	c.errored = make(chan *Result, 1)

	return c
}

func linksToStrings(links []*domain.Link) []string {
	ret := make([]string, len(links))

	for i, l := range links {
		ret[i] = l.Target.String()
	}

	return ret
}

func urlsToStrings(urls []*url.URL) []string {
	ret := make([]string, len(urls))

	for i, u := range urls {
		ret[i] = u.String()
	}

	return ret
}

func strToUrl(s string) *url.URL {
	u, _ := url.Parse(s)
	return u
}

// stringsToUrls converts slices of strings to URLs
func stringsToUrls(strs []string) ([]*url.URL, error) {
	ret := make([]*url.URL, len(strs))

	for i, s := range strs {
		u, err := url.Parse(s)
		if err != nil {
			return nil, err
		}
		ret[i] = u
	}

	return ret, nil
}
