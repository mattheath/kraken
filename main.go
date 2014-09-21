package main

import (
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

	urls := make([]string, 0)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", urls, err
	}

	sel := doc.Find("a")
	if sel == nil {
		return "", urls, nil
	}

	for i, n := range sel.Nodes {
		if n.Type != html.ElementNode || n.DataAtom != atom.A {
			log.Debugf("Node is not an anchor: %v", n.Type)
			continue
		}

		var href string

		for _, a := range n.Attr {
			if a.Key != "href" {
				continue
			}
			href = a.Val
		}

		if href == "" {
			continue
		}

		log.Infof("Node %v: %s", i, href)
		urls = append(urls, href)
	}

	log.Debugf("URLs: %+v", urls)

	return "body", urls, nil
}
