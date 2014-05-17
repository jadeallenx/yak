package yak

import (
	"code.google.com/p/go.net/html"
	"log"
	"net/url"
    "io"
)

type Assettype int

const (
	Script Assettype = iota
	Image
	CSS
	Link
)

type Asset struct {
	Atype Assettype
	Loc   *url.URL
}

type Page struct {
	r        io.Reader  // raw content
	c        *html.Node // parsed (cooked) content
	Parent   *Page
	Loc      *url.URL
	Children []*Page
	Assets   []*Asset
}


type Visit func(n *html.Node, p *Page) error


func Walk(p *Page, v Visit) error {
	if p.c == nil {
		log.Println("WARNING: no html nodes found")
		return nil
	}

	for ch := p.c.FirstChild; ch != nil; ch = ch.NextSibling {
		err := v(ch, p)
		if err != nil {
			return err
		}
		err = walk(ch, p, v)
		if err != nil {
			return err
		}
	}

	return nil

}

func walk(n *html.Node, p *Page, v Visit) error {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err := v(c, p)
		if err != nil {
			return err
		}
		walk(c, p, v)
	}
	return nil
}
