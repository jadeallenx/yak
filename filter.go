package yak

import (
	"code.google.com/p/go.net/html"
	"log"
)

type Visit func(n *html.Node, p *Page) error

func Walk(p *Page, v Visit) *Page {
	if p.c == nil {
		log.Println("WARNING: no html nodes found")
		return p
	}

	for ch := p.c.FirstChild; ch != nil; ch = ch.NextSibling {
		err := v(ch, p)
		if err != nil {
			log.Fatal(err)
		}
		err = walk(ch, p, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	return p

}

func walk(n *html.Node, p *Page, v Visit) error {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err := v(c, p)
		if err != nil {
			return err
			break
		}
		walk(c, p, v)
	}
	return nil
}
