package yak

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
)

type assettype int

const (
	Script assettype = iota
	Image
	CSS
	Link
)

type Asset struct {
	atype assettype
	loc   *url.URL
}

type Page struct {
	r        io.Reader  // raw content
	c        *html.Node // parsed (cooked) content
	parent   *Page
	loc      *url.URL
	children []*Page
	assets   []*Asset
}

func process(q *Queue) error {
	if q.Empty() {
		return nil
	}

	p := q.Dequeue().(*Page)
	err := GetPage(p)
	if err != nil {
		return err
	}

	err = Walk(p, defaultf)
	if err != nil {
		return err
	}

	for _, ch := range p.children {
		q.Enqueue(ch)
	}

	process(q)

	return nil
}

func main() {
	q := NewQueue()
	root := &Page{}
	var err error
	for _, rawurl := range os.Args {
		fmt.Println(rawurl)
		root.loc, err = url.Parse(rawurl)
		if err != nil {
			log.Fatal(err)
		}

		q.Enqueue(root)
	}

	err = process(q)
	if err != nil {
		log.Fatal(err)
	}

	//pretty print p
}
