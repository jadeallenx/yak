package main

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"github.com/mrallen1/yak"
	"log"
	"net/url"
	"os"
)

func main() {
	q := yak.NewQueue()
	root := &yak.Page{}
	var err error
	m := make(map[string]bool)
	for _, rawurl := range os.Args {
		root.Loc, err = url.Parse(rawurl)
		if err != nil {
			log.Fatal(err)
		}

		m[root.Loc.String()] = true
		q.Enqueue(root)
	}

	err = process(q, m)
	if err != nil {
		log.Fatal(err)
	}

	yak.PrettyPrint(root, 1)

	fmt.Println("Done")
}

func process(q *yak.Queue, m map[string]bool) error {
	if q.Empty() {
		return nil
	}

	p := q.Dequeue().(*yak.Page)
	err := yak.GetPage(p)
	if err != nil {
		return err
	}

	err = yak.Walk(p, filterf)
	if err != nil {
		return err
	}

	for _, ch := range p.Children {
		if !m[ch.Loc.String()] {
			m[ch.Loc.String()] = true
			q.Enqueue(ch)
		}
	}

	log.Println("Queue length:", q.Len())

	err = process(q, m)
	if err != nil {
		return err
	}
	return nil
}

func resolveUrl(u *url.URL, p *yak.Page) *url.URL {
	return p.Loc.ResolveReference(u)
}

func addChildPage(parent *yak.Page, u *url.URL) {
	ch := &yak.Page{
		Parent: parent,
		Loc:    u,
	}
	parent.Children = append(parent.Children, ch)
}

func addAsset(parent *yak.Page, u *url.URL, t yak.Assettype) {
	a := &yak.Asset{
		Atype: t,
		Loc:   resolveUrl(u, parent),
	}
	parent.Assets = append(parent.Assets, a)
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

func filterf(n *html.Node, p *yak.Page) error {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a":
			u, err := findAttrKey(n, "href")
			if err != nil {
				// We couldn't parse a URL for this anchor.
				//log.Println("WARNING: no href found for this anchor node", err)
				return nil
			}
			u = resolveUrl(u, p)
			if u.Host == p.Loc.Host {
				// host is the same, child page
				addChildPage(p, u)
			} else {
				// external
				addAsset(p, u, yak.Link)
			}

		case "img":
			u, err := findAttrKey(n, "src")
			if err != nil {
				// We couldn't parse a URL for this image.
				log.Println("WARNING:", err)
				return nil
			}
			addAsset(p, u, yak.Image)
		case "script":
			u, _ := findAttrKey(n, "src")
			if u == nil {
				// This script has no src, so it must be inline
				a := &yak.Asset{Atype: yak.Script}
				p.Assets = append(p.Assets, a)
				return nil
			}
			addAsset(p, u, yak.Script)
		case "link":
			if isCSS(n) {
				u, _ := findAttrKey(n, "href")
				if u == nil {
					// This link has no href, so it must be inline
					a := &yak.Asset{Atype: yak.CSS}
					p.Assets = append(p.Assets, a)
					return nil
				}
				addAsset(p, u, yak.CSS)
			}
		}
	}
	return nil
}
