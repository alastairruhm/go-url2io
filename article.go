package url2io

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/liamylian/jsontime"
)

var json = jsontime.ConfigWithCustomTimeFormat

const (
	articleBasePath = "/article"
)

// ArticleService ...
type ArticleService interface {
	Parse(url string, fields []string) (*http.Response, *ParseResult, error)
}

// ArticleServiceOp impl ArticleService
type ArticleServiceOp struct {
	client *Client
}

// ParseResult normal result of api /article returns
type ParseResult struct {
	Title   string     `json:"title"`
	Content string     `json:"content,omitempty"`
	URL     string     `json:"url"`
	Date    *time.Time `json:"date" time_format:"2006-01-02 15:04:05" time_utc:"false"`
	Text    string     `json:"text,omitempty"`
	Next    string     `json:"next,omitempty"`
}

var _ ArticleService = &ArticleServiceOp{}

// Parse ...
func (s *ArticleServiceOp) Parse(url string, fields []string) (*http.Response, *ParseResult, error) {
	if url == "" {
		return nil, nil, NewArgError("url", "can't be empty")
	}
	for _, f := range fields {
		if f != "next" && f != "text" {
			return nil, nil, NewArgError("fields", "now only support next and text options")
		}
	}
	path := fmt.Sprintf("%s%s", s.client.BaseURL, articleBasePath)

	req, err := http.NewRequest(http.MethodGet, path, nil)

	if err != nil {
		return nil, nil, err
	}

	// construct query string
	q := req.URL.Query()
	q.Add("token", s.client.Token)
	q.Add("url", url)

	if fields != nil {
		q.Add("fields", strings.Join(fields, ","))
	}
	req.URL.RawQuery = q.Encode()

	result := new(ParseResult)
	resp, err := s.Do(req, result)
	if err != nil {
		return nil, nil, err
	}

	return resp, result, nil
}

func (s *ArticleServiceOp) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := s.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	err = CheckResponse(resp)

	if err != nil {
		return nil, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return resp, nil
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	responseError := &ResponseError{}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, responseError)
		if err != nil {
			return err
		}
	}

	return responseError
}

// An ResponseError indicates the response status code isn't 2xx
type ResponseError struct {
	Message string `json:"msg,omitempty"`
	// Error error type
	ErrorType string `json:"error,omitempty"`
	// URL the url which error relates to
	URL string `json:"url"`
	// Code show the http response status code when receive an error response
	Code int `json:"code,omitempty"`
	// show the resource in MIME formation
	ResourceType string `json:"type,omitempty"`
}

func (r *ResponseError) Error() string {
	return fmt.Sprintf("response error code: %d, type: %s, message: %s, url: %s, resource type: %s",
		r.Code,
		r.ErrorType,
		r.Message,
		r.URL,
		r.ResourceType)
}
