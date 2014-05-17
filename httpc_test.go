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
	p := &Page{Loc: u}
	err := GetPage(p)
	if err != nil {
		t.Fatal("Expected nil, got", err)
	}
}

func TestGetPageSecure(t *testing.T) {
	u, _ := url.Parse("https://www.bing.com")
	p := &Page{Loc: u}
	err := GetPage(p)
	if err != nil {
		t.Fatal("Expected nil, got", err)
	}
}

func Test404Page(t *testing.T) {
	u, _ := url.Parse("http://httpstat.us/404")
	p := &Page{Loc: u}
	err := GetPage(p)
	if err != nil {
		t.Fatal("Expected nil, got err")
	}
    if p.c != nil {
        t.Fatal("Expected no parsed content")
    }
}
