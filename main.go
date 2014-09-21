package main

import (
	"errors"
	"net/url"

	html "code.google.com/p/go.net/html"
	atom "code.google.com/p/go.net/html/atom"
	"github.com/PuerkitoBio/goquery"
	log "github.com/cihub/seelog"
)

func main() {

	// Flush logs before exit
	defer log.Flush()

	// Crawl the specified site
	Crawl("http://golang.org/")
}

// Crawl takes a URL and recursively crawls pages
func Crawl(url string) {

	_, urls, err := fetch(url)
	if err != nil {
		log.Errorf("Error:", err)
		return
	}

	log.Infof("URLs found: %+v", urls)
}

// fetch retrieves the page at the specified URL and extracts URLs
func fetch(url string) (string, []string, error) {

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", nil, err
	}

	urls, err := extractLinks(doc)
	if err != nil {
		return "", urls, err
	}

	log.Debugf("URLs: %+v", urls)

	return "", urls, nil
}

// extractLinks from a document
func extractLinks(doc *goquery.Document) ([]string, error) {

	// Blank slice to hold the links on this page
	urls := make([]string, 0)

	// Extract all 'a' elements from the document
	sel := doc.Find("a")
	if sel == nil {
		// Assume zero links on failure
		return nil, nil
	}

	// Range over links, and add them to the list if valid
	for i, n := range sel.Nodes {

		// Validate the node is a link, and extract the target URL
		href, err := validateLink(n)
		if err != nil || href == "" {
			continue
		}

		// Normalise the URL and add if valid
		if uri := normaliseUrl(doc.Url, href); uri != "" {
			log.Debugf("Node %v: %s", i, href)
			urls = append(urls, uri)
		}
	}

	return urls, nil
}

// validateLink is an anchor with a href, and extract normalised url
func validateLink(n *html.Node) (string, error) {
	var href string

	// Confirm this node is an anchor element
	if n == nil || n.Type != html.ElementNode || n.DataAtom != atom.A {
		return href, errors.New("Node is not an anchor")
	}

	// Return the value of the href attr if it exists
	for _, a := range n.Attr {
		if a.Key == "href" && a.Val != "" {
			return a.Val, nil
		}
	}

	return "", errors.New("Node does not contain a href attribute")
}

// normaliseUrl converts relative URLs to absolute URLs
func normaliseUrl(parent *url.URL, urlString string) string {

	// Parse the string into a url.URL
	uri, err := url.Parse(urlString)
	if err != nil {
		log.Debugf("Failed to parse URL: %s", urlString)
		return ""
	}

	// Resolve references to get an absolute URL
	abs := parent.ResolveReference(uri).String()
	log.Debugf("Resolved: %s", abs)

	return abs
}
