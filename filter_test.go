package yak

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	s := `<p>Links:</p><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul>`

	c, _ := html.Parse(strings.NewReader(s))
	fmt.Println(c)
}
