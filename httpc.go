package yak

import (
	"crypto/tls"
	"errors"
	"net/http"
)

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{},
}
var client = &http.Client{Transport: tr}

func GetPage(p *Page) error {
	if p.loc == nil {
		return errors.New("WARNING: location is nil!")
	}
	resp, err := client.Get(p.loc.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	p.r = resp.Body

	return nil
}
