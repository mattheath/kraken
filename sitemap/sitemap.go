package sitemap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/mattheath/kraken/domain"
)

const (
	sitemapTemplateHeader = `<?xml version="1.0" encoding="utf-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
   xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
   xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd">
`

	sitemapTemplateFooter = `</urlset>`

	urlTemplate = `	<url>
		<loc>%s</loc>
		<lastmod>%s</lastmod>
		<changefreq>daily</changefreq>
		<priority>0.8</priority>
	</url>
`
)

type formattedPage struct {
	Url    string   `json:"url"`
	Links  []string `json:"links"`
	Assets []string `json:"assets"`
}

// BuildXMLSitemap builds a standard XML sitemap from a list of pages on a site
func BuildXMLSitemap(pages []*domain.Page) ([]byte, error) {
	var buf bytes.Buffer

	// Add the header with schema definition
	buf.WriteString(sitemapTemplateHeader)

	// Add each page
	for _, p := range pages {
		if p == nil || p.Url == nil {
			continue
		}
		buf.WriteString(fmt.Sprintf(urlTemplate, p.Url.String(), time.Now().Format("2006-01-02")))
	}

	// Append the footer closing tag
	buf.WriteString(sitemapTemplateFooter)

	return buf.Bytes(), nil
}

func BuildJSONSiteStructure(target *url.URL, pages []*domain.Page) ([]byte, error) {

	ret := map[string]interface{}{
		"target": target.String(),
	}

	ps := []*formattedPage{}
	for _, p := range pages {
		fp := &formattedPage{
			Url: p.Url.String(),
		}

		fp.Links = make([]string, len(p.Links))
		for i, l := range p.Links {
			fp.Links[i] = l.Target.String()
		}

		fp.Assets = make([]string, len(p.Assets))
		for i, a := range p.Assets {
			fp.Assets[i] = a.String()
		}

		ps = append(ps, fp)
	}
	ret["pages"] = ps

	return json.Marshal(ret)
}
