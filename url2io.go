package url2io

import (
	"net/http"
	"net/url"
)

const (
	// Version reveals SDK version
	Version        = "1.0.0"
	defaultBaseURL = "http://api.url2io.com"
	userAgent      = "url2io/" + Version
	mediaType      = "application/json"
)

// Client manages communication with url2io API.
type Client struct {
	client *http.Client
	// BaseURL string
	BaseURL   *url.URL
	Token     string
	UserAgent string

	Article ArticleService
}

// NewClient returns a new client of url2io
func NewClient(token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, Token: token, BaseURL: baseURL, UserAgent: userAgent}
	c.Article = &ArticleServiceOp{client: c}
	return c
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
