package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	html "code.google.com/p/go.net/html"
	atom "code.google.com/p/go.net/html/atom"
)

func TestExtractValidHrefSuccess(t *testing.T) {
	f := &HttpFetcher{}

	successCases := map[string]*html.Node{
		"https://example.com": &html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.A,
			Attr: []html.Attribute{
				html.Attribute{
					Key: "href",
					Val: "https://example.com",
				},
			},
		},
		"/doc/": &html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.A,
			Attr: []html.Attribute{
				html.Attribute{
					Key: "title",
					Val: "Many links, such wow",
				},
				html.Attribute{
					Key: "href",
					Val: "/doc/",
				},
			},
		},
	}

	for expected, tc := range successCases {
		res, err := f.extractValidHref(tc)
		assert.Nil(t, err)
		assert.Equal(t, res, expected)
	}
}

func TestExtractValidHrefFailure(t *testing.T) {
	f := &HttpFetcher{}

	// Also test a number of failure cases
	failureCases := map[error]*html.Node{

		InvalidNode: &html.Node{},
		InvalidNode: &html.Node{
			Type: html.DocumentNode,
		},
		InvalidNode: &html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.Div,
		},

		InvalidNodeAttributeMissing: &html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.A,
			Attr:     []html.Attribute{},
		},
		InvalidNodeAttributeMissing: &html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.A,
			Attr: []html.Attribute{
				html.Attribute{
					Key: "title",
					Val: "Many links, such wow",
				},
			},
		},
	}

	for expected, tc := range failureCases {
		res, err := f.extractValidHref(tc)

		assert.NotNil(t, err)
		assert.Equal(t, res, "")

		assert.Equal(t, err, expected)
	}
}

func TestNormaliseUrl(t *testing.T) {
	f := &HttpFetcher{}
	parent := &url.URL{
		Scheme: "http",
		Host:   "example.com",
		Path:   "/",
	}

	testCases := map[string]string{
		"/boop/":                                      "http://example.com/boop/",
		"/shoop":                                      "http://example.com/shoop",
		"dawhoop":                                     "http://example.com/dawhoop",
		"http://notrelative.com/r/boop/":              "http://notrelative.com/r/boop/",
		"http://tinyhamsterseatingtinyburritos":       "http://tinyhamsterseatingtinyburritos",
		"https://www.youtube.com/watch?v=JOCtdw9FG-s": "https://www.youtube.com/watch?v=JOCtdw9FG-s",
	}

	for tc, expected := range testCases {
		result := f.normaliseUrl(parent, tc)
		assert.Equal(t, result, expected)
	}

}
