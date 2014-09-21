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
	Crawl(*target, 4, fetcher)
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
