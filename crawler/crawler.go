package crawler

import (
	"net/url"

	log "github.com/cihub/seelog"
)

// Crawler coordinated crawling a site, and stores completed results
type Crawler struct {
	// Store our results
	Pages map[string]*Page
	Links map[string]*Link

	// completed channel is an inbound queue of completed requests
	// for processing by the main crawler goroutine
	completed chan *Result

	// skipped tracks pages we have skipped
	skipped chan *Result

	// errored tracks pages which errored, which we may then
	// choose to reattempt
	errored chan *Result

	// requestsInFlight tracks how many of requests are outstanding
	requestsInFlight int

	// totalRequests tracks the number of requests we have made
	totalRequests int
}

type Result struct {
	Url   *url.URL
	Depth int
	Page  *Page
	Error error
}

// Work is our main event loop, coordinating request processing
// This is single threaded and is the only thread that writes into
// our internal maps, so we don't require coordination or locking
// (maps are not threadsafe)
func (c *Crawler) Work(target string, depth int, fetcher Fetcher) {

	// Convert our target to a URL
	uri, err := url.Parse(target)
	if err != nil {
		log.Errorf("Could not parse target '%s'", target)
		return
	}

	// Initialise channels to track requests
	c.completed = make(chan *Result)
	c.skipped = make(chan *Result)
	c.errored = make(chan *Result)

	// Initialise results containers
	c.Pages = make(map[string]*Page)
	c.Links = make(map[string]*Link)

	// Get our first page & track this
	go c.crawl(uri, depth, fetcher)
	c.requestsInFlight++

	// Event loop
	for {
		select {
		case r := <-c.skipped:
			log.Debugf("Page skipped for %s", r.Url)
			c.totalRequests--
		case r := <-c.errored:
			log.Debugf("Page errored for %s: %v", r.Url, r.Error)
		case r := <-c.completed:
			log.Debugf("Page complete for %s", r.Url)

			if r.Page == nil {
				break
			}

			for _, l := range r.Page.Links {
				go c.crawl(l.Target, r.Depth-1, fetcher)
				c.requestsInFlight++
				c.totalRequests++
			}
			log.Debugf("Fired %v new requests, %v currently in flight", len(r.Page.Links), c.requestsInFlight)

			c.Pages[r.Url.String()] = r.Page

		}

		// Decrement outstanding requests & and abort if complete
		c.requestsInFlight--
		if c.requestsInFlight == 0 {
			log.Debugf("Complete")
			return
		}
	}
}

// crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func (c *Crawler) crawl(source *url.URL, depth int, fetcher Fetcher) {

	// The result of our crawl
	res := &Result{
		Depth: depth,
		Url:   source,
	}

	// Skip pages if we are at our maximum depth
	if depth <= 0 {
		log.Debugf("Skipping %s as at 0 depth", source.String())
		c.skipped <- res
		return
	}

	// Crawl the page, using our fetcher
	_, urls, err := fetcher.Fetch(source.String())
	if err != nil {
		res.Error = err
		c.errored <- res
		return
	}

	log.Infof("%v URLs found at %s", len(urls), source.String())

	links := make([]*Link, 0)
	for _, u := range urls {
		if nu, err := url.Parse(u); err == nil {
			links = append(links, &Link{
				Source: source,
				Target: nu,
			})
		}
	}

	// Store this page and links into the result
	res.Page = &Page{
		Url:   source,
		Links: links,
	}

	// 	// Mark this page as complete
	c.completed <- res
}

func (c *Crawler) TotalRequests() int {
	return c.totalRequests
}
