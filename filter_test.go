package yak

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"log"
	"net/url"
	"strings"
	"testing"
)

func TestWalk(t *testing.T) {
	s := `<link type="text/css" rel="stylesheet" href="style.css"><p>Links:</p><img src="qux.jpg"><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul><a href="#">lala</a>`

	p := &Page{
		r:        strings.NewReader(s),
		Parent:   nil,
		Children: nil,
		Assets:   nil,
	}

	p.c, _ = html.Parse(p.r)
	p.Loc, _ = url.Parse("https://test.com/")

	err := Walk(p, testf)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.Children) != 3 {
		t.Fatalf("Expected child length of 3, got %v", len(p.Children))
	}

	if len(p.Assets) != 2 {
		t.Fatalf("Expected assets length of 1, got %v", len(p.Assets))
	}

}

func TestBadWalkInput(t *testing.T) {
	s := `<a>Degenerate Input</a><img><script>Fail whale</script>`

	p := &Page{
		r:        strings.NewReader(s),
		Parent:   nil,
		Children: nil,
		Assets:   nil,
	}

	p.c, _ = html.Parse(p.r)
	p.Loc, _ = url.Parse("https://test.com/")

	err := Walk(p, testf)
	if err != nil {
		t.Fatal(err)
	}

}

func testf(n *html.Node, p *Page) error {
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
			if u.Host == p.Loc.Host {
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
				a := &Asset{Atype: Script}
				p.Assets = append(p.Assets, a)
				return nil
			}
			addAsset(p, u, Script)
		case "link":
			if isCSS(n) {
				u, _ := findAttrKey(n, "href")
				if u == nil {
					// This link has no href, so it must be inline
					a := &Asset{Atype: CSS}
					p.Assets = append(p.Assets, a)
					return nil
				}
				addAsset(p, u, CSS)
			}
		}
	}
	return nil
}

func resolveUrl(u *url.URL, p *Page) *url.URL {
	return p.Loc.ResolveReference(u)
}

func addChildPage(parent *Page, u *url.URL) {
	ch := &Page{
		Parent: parent,
		Loc:    u,
	}
	parent.Children = append(parent.Children, ch)
}

func addAsset(parent *Page, u *url.URL, t Assettype) {
	a := &Asset{
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
