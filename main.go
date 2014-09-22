package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	log "github.com/cihub/seelog"

	"github.com/mattheath/kraken/crawler"
	"github.com/mattheath/kraken/sitemap"
)

var (
	flagSet = flag.NewFlagSet("kraken", flag.ExitOnError)

	target         = flagSet.String("target", "", "target URL to crawl")
	depth          = flagSet.Int("depth", 4, "depth of pages to crawl")
	verboseLogging = flagSet.Bool("v", false, "enable verbose logging")
	outputDir      = flagSet.String("o", "", "directory to output to")
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

	// Directory to save output files
	out := *outputDir
	if out == "" {
		out, err = os.Getwd()
		if err != nil {
			log.Criticalf("Failed to get current working directory: %v", err)
		}
	}

	// Use a HTTP based fetcher
	fetcher := &HttpFetcher{}

	// Fire!
	log.Infof("Unleashing the Kraken at %s", *target)

	// Crawl the specified site
	c := &crawler.Crawler{}
	c.Work(targetUrl, *depth, fetcher)

	// Success
	log.Infof("%v pages found, %v requests attempted", len(c.Pages), c.TotalRequests())

	writeSitemaps(out, c)
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

func writeSitemaps(outdir string, c *crawler.Crawler) error {

	// Build sitemap and write to output file
	xmlout := fmt.Sprintf("%s/%s-sitemap.xml", outdir, c.Target().Host)
	xmlSitemap, err := sitemap.BuildXMLSitemap(c.AllPages())
	if err != nil {
		log.Criticalf("Failed to generate sitemap to %s", xmlout)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(xmlout, xmlSitemap, 0644); err != nil {
		log.Criticalf("Failed to write sitemap to %s", xmlout)
		os.Exit(1)
	}
	log.Infof("Wrote XML sitemap to %s", xmlout)

	// Build JSON site description
	siteout := fmt.Sprintf("%s/%s-sitemap.json", outdir, c.Target().Host)

	b, err := sitemap.BuildJSONSiteStructure(c.Target(), c.AllPages())

	if err := ioutil.WriteFile(siteout, b, 0644); err != nil {
		log.Criticalf("Failed to write sitemap to %s", siteout)
		os.Exit(1)
	}
	log.Infof("Wrote JSON sitemap to %s", siteout)

	return nil
}
