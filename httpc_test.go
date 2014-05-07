package yak

import (
	"net/url"
	"testing"
)

func TestGetPageNil(t *testing.T) {
	p := &Page{}
	err := GetPage(p)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestGetPage(t *testing.T) {
	u, _ := url.Parse("http://www.example.com")
	p := &Page{loc: u}
	err := GetPage(p)
	if err != nil {
		t.Fatal("Expected non-nil, got", err)
	}
}

func TestGetPageSecure(t *testing.T) {
	u, _ := url.Parse("https://www.bing.com")
	p := &Page{loc: u}
	err := GetPage(p)
	if err != nil {
		t.Fatal("Expected non-nil, got", err)
	}
}

