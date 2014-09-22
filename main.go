package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/cihub/seelog"

	"github.com/mattheath/kraken/crawler"
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
	c := &crawler.Crawler{}
	c.Work(*target, *depth, fetcher)

	log.Debugf("We're done!")
}
