package sidelines

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type Endpoint struct {
	Method   string
	Url      *url.URL
	ProxyUrl *url.URL
	PostData url.Values
	Client   *http.Client
}

type Response struct {
	StatusCode int
	Body       []byte
}

func (endpoint *Endpoint) Request() (*Response, error) {
	if endpoint.Client == nil {
		if endpoint.ProxyUrl == nil {
			endpoint.Client = http.DefaultClient
		} else {
			endpoint.Client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(endpoint.ProxyUrl)}}
		}
	}

	var rawResp *http.Response
	var err error

	if endpoint.Method == "GET" {
		rawResp, err = endpoint.Client.Get(endpoint.Url.String())
	} else if endpoint.Method == "POST" {
		rawResp, err = endpoint.Client.PostForm(endpoint.Url.String(), endpoint.PostData)
	}
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(rawResp.Body)
	rawResp.Body.Close()

	if err != nil {
		return nil, err
	}

	resp := &Response{
		StatusCode: rawResp.StatusCode,
		Body:       body,
	}

	return resp, nil
}
