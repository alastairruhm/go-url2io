package url2io

import (
	"net/http"
	"time"
)

const (
	articleBasePath = "/article"
)

type ArticleService interface {
	Parse(url string, fields []string) (*ParseResult, *http.Response, error)
}

type ArticleServiceOp struct {
	client *Client
}

type ParseResult struct {
	Title   string    `json:"title"`
	Content string    `json:"content,omitempty"`
	Url     string    `json:"url"`
	Date    time.Time `json:"date"`
	Text    string    `json:"text,omitempty"`
	Next    string    `json:"next,omitempty"`
}

var _ ArticleService = &ArticleServiceOp{}

func (s *ArticleServiceOp) Parse(url string, fields []string) (*ParseResult, *http.Response, error) {
	return nil, nil, nil
}
