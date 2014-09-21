package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/cihub/seelog"
)

var (
	flagSet = flag.NewFlagSet("kraken", flag.ExitOnError)

	target = flagSet.String("target", "", "target URL to crawl")
	depth  = flagSet.Int("depth", 2, "depth of pages to crawl")
)

func main() {
	// Process flags
	flagSet.Parse(os.Args[1:])

	// Flush logs before exit
	defer log.Flush()

	// Do we have a target?
	if *target == "" {
		fmt.Println("Please specify a target domain, eg. kraken -target=\"http://example.com\"")
		os.Exit(1)
	}
	log.Infof("Unleashing the Kraken at %s", *target)

	// Use a HTTP based fetcher
	fetcher := &HttpFetcher{}

	// Crawl the specified site
	done := make(chan bool, 1)
	Crawl(*target, *depth, fetcher, done)
	<-done

	log.Debugf("We're done!")
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, pageDone chan bool) {

	if depth <= 0 {
		log.Debugf("Skipping %s as at 0 depth", url)
		pageDone <- true
		return
	}

	_, urls, err := fetcher.Fetch(url)
	if err != nil {
		log.Errorf("Error:", err)
		pageDone <- true
		return
	}

	log.Infof("%v URLs found at %s", len(urls), url)

	// Track children
	done := make(chan bool)
	count := 0

	for _, u := range urls {
		log.Debugf("Firing crawler at %s, depth %v", u, depth-1)
		count++
		go Crawl(u, depth-1, fetcher, done)
	}

	for ; count > 0; count-- {
		log.Debugf("waiting on done chan")
		<-done
	}

	log.Debugf("Page complete: %s", url)

	// Mark this page as complete
	pageDone <- true
}
