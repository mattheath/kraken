package crawler

import (
	"net/url"
)

type Fetcher interface {
	// Fetch returns a slice of URLs found on the target page
	// along with a slice of assets.
	Fetch(target *url.URL) (urls []*url.URL, assets []*url.URL, err error)
}
