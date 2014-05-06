package yak

import (
	"code.google.com/p/go.net/html"
	"io"
	"net/url"
)

type assettype int

const (
	Script assettype = iota
	Image
	CSS
)

type Asset struct {
	atype assettype
	loc   *url.URL
}

type Page struct {
	r        *io.Reader // raw content
	c        *html.Node // parsed (cooked) content
	parent   *Page
	loc      *url.URL
	children []*Page
	assets   []Asset
}
