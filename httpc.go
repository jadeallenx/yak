package yak

import (
	"code.google.com/p/go.net/html"
	"crypto/tls"
	"errors"
	"log"
	"net/http"
)

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{},
}
var client = &http.Client{Transport: tr}

func GetPage(p *Page) error {
	if p.Loc == nil {
		return errors.New("WARNING: location is nil!")
	}
	resp, err := client.Get(p.Loc.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Got *something* in the 200 range. Let's be optimistic and (try to) parse it.
	if resp.StatusCode > 199 && resp.StatusCode < 300 {
		p.r = resp.Body
		p.c, err = html.Parse(p.r)
		if err != nil {
            log.Println(err)
			return nil
		}
	} else {
        // Got something other than a 200, so just error
		log.Printf("WARNING: %v returned status code: %v", p.Loc.String(), resp.StatusCode)
	}

	return nil
}
