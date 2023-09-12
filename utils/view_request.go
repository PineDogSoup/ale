package utils

import (
	"net/http"
	"net/url"
)

type viewRequest struct {
}

func (a *viewRequest) HTTPRequest(ep url.URL) *http.Request {
	params := ep.Query()
	ep.RawQuery = params.Encode()
	req, _ := http.NewRequest("GET", ep.String(), nil)
	return req
}
