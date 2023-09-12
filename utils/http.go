package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	Endpoints         []string
	Version           string
	TimeoutPerRequest time.Duration
}

type HttpClient struct {
	endpoints []string
	Version   string
	Client    *http.Client
}

func NewHttpClient(cfg Config) *HttpClient {
	c := &HttpClient{
		Version: cfg.Version,
		Client: &http.Client{
			Timeout: cfg.TimeoutPerRequest,
		},
	}

	return c
}

func (c *HttpClient) Do(method, url string, params map[string]interface{}) (*http.Response, []byte, error) {
	var headerKey string
	var reader io.Reader
	var body []byte
	var resp *http.Response

	switch method {
	case "POST":
		if params != nil {
			bytesData, _ := json.Marshal(params)
			reader = strings.NewReader(string(bytesData))
		}
		headerKey = "Content-Type"
	case "GET":
		if params != nil {
			strParams := Map2UrlParams(params)
			url += "?" + strParams
		}
		headerKey = "Accept"
	default:
		return resp, body, errors.New("unexpected http method")
	}

	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return resp, body, err
	}

	if c.Version != "" {
		request.Header.Set(headerKey, "application/json;v="+c.Version)
	} else {
		request.Header.Set(headerKey, "application/json")
	}

	if resp, err = c.Client.Do(request); err != nil {
		return resp, body, err
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if body, err = io.ReadAll(resp.Body); err != nil {
		body = []byte(fmt.Sprintf("(could not fetch response reader for error: %s)", err))
	}

	return resp, body, err
}

type httpAction interface {
	HttpRequest(url.URL) *http.Request
}

type callAction struct {
}

func (c callAction) HttpRequest(url url.URL) *http.Request {
	params := url.Query()
	url.RawQuery = params.Encode()
	req, _ := http.NewRequest("GET", url.String(), nil)
	return req
}

type sendAction struct {
	Data   []byte
	Params map[string]interface{}
}

func (s *sendAction) HttpRequest(u url.URL) *http.Request {
	bytesData, _ := json.Marshal(s.Params)
	body := strings.NewReader(string(bytesData))
	req, _ := http.NewRequest("POST", u.String(), body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}
