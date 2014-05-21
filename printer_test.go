package yak

import (
	"code.google.com/p/go.net/html"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func chooseType(r *rand.Rand) (Assettype, string) {
	switch r.Intn(10) {
	case 0:
		return Link, ".link"
	case 1:
		return CSS, ".css"
	case 2:
		return Script, ".js"
	default:
		return Image, ".jpg"
	}
}

func generate_page(r *rand.Rand, baseUrl string, baseOnly bool) *Page {

	s := `<a href="foo">Foo</a>`

	p := &Page{
		r:        strings.NewReader(s),
		Parent:   nil,
		Children: nil,
		Assets:   nil,
	}

	p.c, _ = html.Parse(p.r)

    if baseOnly {
	    p.Loc, _ = url.Parse(baseUrl)
    } else {
	    p.Loc, _ = url.Parse(baseUrl + "/" + strconv.FormatInt(r.Int63n(10000), 36))
    }

	for i := r.Intn(5); i != 0; i-- {
		at, ext := chooseType(r)
		u, _ := url.Parse(strconv.FormatInt(r.Int63n(10000), 36) + ext)
		a := &Asset{
			Atype: at,
			Loc:   p.Loc.ResolveReference(u),
		}
		p.Assets = append(p.Assets, a)
	}
	return p
}

func generate_children(r *rand.Rand, baseUrl string) []*Page {
	var children []*Page
	for i := r.Intn(5); i != 0; i-- {
		ch := generate_page(r, baseUrl, false)
		children = append(children, ch)
	}
	return children
}

func TestPrinter(t *testing.T) {
	// crap entropy but doesn't need to be cryptographically strong
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// build a tree of pages
	parent := generate_page(r, "http://test.com", true)

	for i := r.Intn(5) + 1; i != 0; i-- {
		ch := generate_page(r, parent.Loc.String(), false)
		ch.Children = generate_children(r, ch.Loc.String())
		parent.Children = append(parent.Children, ch)
	}

	PrettyPrint(parent, 1)
}
