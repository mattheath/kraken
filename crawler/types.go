package crawler

import (
	"net/url"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (urls []*url.URL, assets []*url.URL, err error)
}

type Page struct {
	Url    *url.URL
	Links  []*Link
	Assets []*url.URL
}

type Link struct {
	Source *url.URL
	Target *url.URL
}
