package yak

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"log"
	"net/url"
)

type Visit func(n *html.Node, p *Page) error

func resolveUrl(u *url.URL, p *Page) *url.URL {
	return p.loc.ResolveReference(u)
}

func addChildPage(parent *Page, u *url.URL) {
	ch := &Page{
		parent: parent,
		loc:    u,
	}
	parent.children = append(parent.children, ch)
}

func addAsset(parent *Page, u *url.URL, t assettype) {
	a := &Asset{
		atype: t,
		loc:   resolveUrl(u, parent),
	}
	parent.assets = append(parent.assets, a)
}

func findAttrKey(n *html.Node, key string) (*url.URL, error) {
	for _, a := range n.Attr {
		if a.Key == key {
			u, err := url.Parse(a.Val)
			if err != nil {
				return nil, err
			}
			return u, nil
		}
	}
	return nil, fmt.Errorf("key '%s' not found", key)
}

func isCSS(n *html.Node) bool {
	for _, a := range n.Attr {
		return a.Key == "type" && a.Val == "text/css"
	}
	return false
}

func defaultf(n *html.Node, p *Page) error {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a":
			u, err := findAttrKey(n, "href")
			if err != nil {
				// We couldn't parse a URL for this anchor.
				log.Println("WARNING:", err)
				return nil
			}
			u = resolveUrl(u, p)
			if u.Host == p.loc.Host {
				// host is the same, child page
				addChildPage(p, u)
			} else {
				// external
				addAsset(p, u, Link)
			}

		case "img":
			u, err := findAttrKey(n, "src")
			if err != nil {
				// We couldn't parse a URL for this image.
				log.Println("WARNING:", err)
				return nil
			}
			addAsset(p, u, Image)
		case "script":
			u, _ := findAttrKey(n, "src")
			if u == nil {
				// This script has no src, so it must be inline
				a := &Asset{atype: Script}
				p.assets = append(p.assets, a)
				return nil
			}
			addAsset(p, u, Script)
		case "link":
			if isCSS(n) {
				u, _ := findAttrKey(n, "href")
				if u == nil {
					// This link has no href, so it must be inline
					a := &Asset{atype: CSS}
					p.assets = append(p.assets, a)
					return nil
				}
				addAsset(p, u, CSS)
			}
		}
	}
	return nil
}

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
