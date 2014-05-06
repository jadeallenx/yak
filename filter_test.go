package yak

import (
	"code.google.com/p/go.net/html"
	"net/url"
	"strings"
	"testing"
)

func TestWalk(t *testing.T) {
	s := `<link type="text/css" rel="stylesheet" href="style.css"><p>Links:</p><img src="qux.jpg"><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul><a href="#">lala</a>`

	p := &Page{
		r:        strings.NewReader(s),
		parent:   nil,
		children: nil,
		assets:   nil,
	}

	p.c, _ = html.Parse(p.r)
	p.loc, _ = url.Parse("https://test.com/")

	err := Walk(p, defaultf)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.children) != 3 {
		t.Fatalf("Expected child length of 3, got %v", len(p.children))
	}

	if len(p.assets) != 2 {
		t.Fatalf("Expected assets length of 1, got %v", len(p.assets))
	}

}

func TestBadWalkInput(t *testing.T) {
	s := `<a>Degenerate Input</a><img><script>Fail whale</script>`

	p := &Page{
		r:        strings.NewReader(s),
		parent:   nil,
		children: nil,
		assets:   nil,
	}

	p.c, _ = html.Parse(p.r)
	p.loc, _ = url.Parse("https://test.com/")

	err := Walk(p, defaultf)
	if err != nil {
		t.Fatal(err)
	}

}
