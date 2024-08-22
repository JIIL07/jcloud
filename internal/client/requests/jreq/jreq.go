package jreq

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/lib/params"
	"io"
	"net/http"
)

func Get(p *params.Params) (*http.Response, error) {
	reqType := p.String("type")
	url := p.String("url")

	req, err := http.NewRequest(reqType, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return resp, nil
}

func Post(p *params.Params) (*http.Response, error) {
	reqType := p.String("type")
	url := p.String("url")
	body := p.Get("body")
	req, err := http.NewRequest(reqType, url, body.(io.Reader))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return resp, nil
}

func Delete(p *params.Params) (*http.Response, error) {
	reqType := p.String("type")
	url := p.String("url")
	req, err := http.NewRequest(reqType, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return resp, nil
}
