package crawler

import (
	log "github.com/cihub/seelog"
)

type Crawler struct {
	// Store our results
	Pages map[string]*Page
	Links map[string]*Link

	// completed channel is an inbound queue of completed requests
	// for processing by the main crawler goroutine
	completed chan *Result

	skipped chan *Result

	errored chan *Result

	// requestsInFlight tracks how many of requests are outstanding
	requestsInFlight int
}

type Result struct {
	Url   string
	Depth int
	Page  *Page
	Error error
}

func (c *Crawler) Work(url string, depth int, fetcher Fetcher) {

	// Initialise a channel to track completed pages
	c.completed = make(chan *Result)

	// Track skipped pages (eg. off site, beyond depth)
	c.skipped = make(chan *Result)

	// Get our first page & track this
	go c.crawl(url, depth, fetcher)
	c.requestsInFlight++

	// Event loop
	for {
		select {
		case r := <-c.skipped:
			log.Debugf("Page skipped for %s", r.Url)
		case r := <-c.errored:
			log.Debugf("Page errored for %s: %v", r.Url, r.Error)
		case r := <-c.completed:
			log.Debugf("Page complete for %s", r.Url)
		}

		// Decrement outstanding requests & and abort if complete
		c.requestsInFlight--
		if c.requestsInFlight == 0 {
			log.Debugf("Complete")
			return
		}

	}

}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func (c *Crawler) crawl(url string, depth int, fetcher Fetcher) {

	res := &Result{
		Depth: depth,
		Url:   url,
	}

	if depth <= 0 {
		log.Debugf("Skipping %s as at 0 depth", url)
		c.skipped <- res
		return
	}

	_, urls, err := fetcher.Fetch(url)
	if err != nil {
		res.Error = err
		c.errored <- res
		return
	}

	log.Infof("%v URLs found at %s", len(urls), url)

	// 	for _, u := range urls {
	// 		log.Debugf("Firing crawler at %s, depth %v", u, depth-1)
	// 		count++
	// 		go Crawl(u, depth-1, fetcher, done)
	// 	}

	// 	for ; count > 0; count-- {
	// 		log.Debugf("waiting on done chan")
	// 		<-done
	// 	}

	// 	log.Debugf("Page complete: %s", url)

	// 	// Mark this page as complete
	c.completed <- res
}
