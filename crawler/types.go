package crawler

import (
	"net/url"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type Page struct {
	Url    *url.URL
	Links  []*Link
	Assets []string
}

type Link struct {
	Source *url.URL
	Target *url.URL
}
