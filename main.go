package main

import (
	log "github.com/cihub/seelog"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

func main() {

	// Flush logs before exit
	defer log.Flush()

	fetcher := &HttpFetcher{}

	// Crawl the specified site
	Crawl("http://golang.org/", 4, fetcher)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {

	_, urls, err := fetcher.Fetch(url)
	if err != nil {
		log.Errorf("Error:", err)
		return
	}

	log.Infof("URLs found: %+v", urls)
}
