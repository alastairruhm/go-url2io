package url2io

import (
	"net/http"
)

const (
	defaultBaseURL = "http://api.url2io.com"
)

type Client struct {
	client *http.Client
	// BaseURL *url.URL
	BaseURL   string
	Token     string
	UserAgent string

	Article ArticleService
}
