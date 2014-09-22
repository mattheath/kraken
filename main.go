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

	target         = flagSet.String("target", "", "target URL to crawl")
	depth          = flagSet.Int("depth", 4, "depth of pages to crawl")
	verboseLogging = flagSet.Bool("v", false, "enable verbose logging")
)

func main() {
	// Process flags
	flagSet.Parse(os.Args[1:])

	// Flush logs before exit
	setLogger(*verboseLogging)
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
	log.Infof("%v pages found, %v requests attempted", len(c.Pages), c.TotalRequests())
}

func setLogger(verbose bool) {

	var logLevel string
	if verbose {
		logLevel = "debug"
	} else {
		logLevel = "info"
	}

	logConfig := `
<seelog>
    <outputs>
        <filter levels="%s">
            <console />
        </filter>
    </outputs>
</seelog>`

	logger, _ := log.LoggerFromConfigAsBytes([]byte(fmt.Sprintf(logConfig, logLevel)))
	log.UseLogger(logger)
}
