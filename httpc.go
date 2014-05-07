package yak

import (
	"code.google.com/p/go.net/html"
	"crypto/tls"
	"errors"
	"fmt"
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

	// Got *something* in the 200 range. Let's be optimistic and (try to) parse it.
	if resp.StatusCode > 199 && resp.StatusCode < 300 {
		p.r = resp.Body
		p.c, err = html.Parse(p.r)
		if err != nil {
			return err
		}
	} else {
        // Got something other than a 200, so just error
		return fmt.Errorf("Status code: %v", resp.StatusCode)
	}

	return nil
}
