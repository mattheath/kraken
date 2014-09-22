package crawler

import (
	"net/url"

	log "github.com/cihub/seelog"

	"github.com/mattheath/kraken/domain"
)

// Crawler coordinated crawling a site, and stores completed results
type Crawler struct {
	// Store our results
	Pages map[string]*domain.Page
	Links map[string]*domain.Link

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

	// target stores our original target for comparisons
	target *url.URL
}

// NewCrawler initialises and returns a new Crawler
func NewCrawler() *Crawler {

	// Initialise new Crawler
	c := &Crawler{
		// Initialise channels to track requests
		completed: make(chan *Result),
		skipped:   make(chan *Result),
		errored:   make(chan *Result),

		// Initialise results containers
		Pages: make(map[string]*domain.Page),
		Links: make(map[string]*domain.Link),
	}

	return c
}

func (c *Crawler) AllPages() []*domain.Page {
	ret := make([]*domain.Page, len(c.Pages))

	// Iterate over pages and assign to slice
	i := 0
	for _, p := range c.Pages {
		ret[i] = p
		i++
	}

	return ret
}

func (c *Crawler) Target() *url.URL {
	return c.target
}

func (c *Crawler) TotalRequests() int {
	return c.totalRequests
}

// Result represents the result of a crawl request
type Result struct {
	Url   *url.URL
	Depth int
	Page  *domain.Page
	Error error
}

// Work is our main event loop, coordinating request processing
// This is single threaded and is the only thread that writes into
// our internal maps, so we don't require coordination or locking
// (maps are not threadsafe)
func (c *Crawler) Work(target *url.URL, depth int, fetcher Fetcher) {

	// Store our target to a URL
	c.target = target

	// Get our first page & track this
	go c.crawl(c.target, depth, fetcher)
	c.requestsInFlight++
	c.totalRequests++

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

			// Process each link
			for _, l := range r.Page.Links {

				// Skip page if not on our target domain
				if l.Target.Host != c.target.Host {
					// log.Debugf("Skipping %s as not on target domain", source.String())
					continue
				}

				// Check if we've already hit this page
				if _, exists := c.Pages[l.Target.String()]; exists {
					// log.Debugf("Skipping %s as already processed", l.Target.String())
					continue
				}

				log.Debugf("Triggering crawl of %s from %s", l.Target.String(), r.Url.String())
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
	urls, assets, err := fetcher.Fetch(source)
	if err != nil {
		res.Error = err
		c.errored <- res
		return
	}

	log.Infof("%v URLs found at %s", len(urls), source.String())

	links := make([]*domain.Link, 0)
	for _, u := range urls {
		links = append(links, &domain.Link{
			Source: source,
			Target: u,
		})
	}

	// Store this page and links into the result
	res.Page = &domain.Page{
		Url:    source,
		Links:  links,
		Assets: assets,
	}

	// 	// Mark this page as complete
	c.completed <- res
}
