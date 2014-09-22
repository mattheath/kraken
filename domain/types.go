package domain

import (
	"net/url"
)

type Page struct {
	Url    *url.URL
	Links  []*Link
	Assets []*url.URL
}

type Link struct {
	Source *url.URL
	Target *url.URL
}
