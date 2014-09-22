package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	log "github.com/cihub/seelog"

	"github.com/mattheath/kraken/crawler"
)

var (
	flagSet = flag.NewFlagSet("kraken", flag.ExitOnError)

	target         = flagSet.String("target", "", "target URL to crawl")
	depth          = flagSet.Int("depth", 4, "depth of pages to crawl")
	verboseLogging = flagSet.Bool("v", false, "enable verbose logging")
	outputFile     = flagSet.String("o", "", "output sitemap to file")
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
	targetUrl, err := url.Parse(*target)
	if err != nil {
		fmt.Println("Could not parse target url '%s' - %v", *target, err)
		os.Exit(1)
	}

	// Save output file
	out := *outputFile
	if out == "" {
		pwd, err := os.Getwd()
		if err != nil {
			log.Criticalf("Failed to get current working directory: %v", err)
		}
		out = fmt.Sprintf("%s/%s-sitemap.xml", pwd, targetUrl.Host)
	}

	// Use a HTTP based fetcher
	fetcher := &HttpFetcher{}

	// Fire!
	log.Infof("Unleashing the Kraken at %s", *target)

	// Crawl the specified site
	c := &crawler.Crawler{}
	c.Work(targetUrl, *depth, fetcher)

	log.Debugf("We're done!")
	log.Infof("%v pages found, %v requests attempted", len(c.Pages), c.TotalRequests())

	log.Infof("Outputting site map to %s", out)
	// output map...
}

// setLogger initialises the logger with the desired verbosity level
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
