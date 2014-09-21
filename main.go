package main

import (
	"flag"
	"fmt"
	"os"

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

	var target string

	// Process flags
	flag.StringVar(&target, "target", "", "Target URL to crawl")
	flag.Parse()

	// Good to go?
	if target == "" {
		fmt.Println("Please specify a target domain, eg. kraken -target=\"http://example.com\"")
		os.Exit(1)
	}
	log.Infof("Unleashing the Kraken at %s", target)

	// Use a HTTP based fetcher
	fetcher := &HttpFetcher{}

	// Crawl the specified site
	Crawl(target, 4, fetcher)
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
