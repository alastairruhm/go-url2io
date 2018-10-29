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

type ArticleService interface {
	Parse(url string, fields []string) (*http.Response, *ParseResult, *ErrorResponse, error)
}

type ArticleServiceOp struct {
	client *Client
}

type ParseResult struct {
	Title   string     `json:"title"`
	Content string     `json:"content,omitempty"`
	URL     string     `json:"url"`
	Date    *time.Time `json:"date" time_format:"2006-01-02 15:04:05" time_utc:"false"`
	Text    string     `json:"text,omitempty"`
	Next    string     `json:"next,omitempty"`
}

var _ ArticleService = &ArticleServiceOp{}

func (s *ArticleServiceOp) Parse(url string, fields []string) (*http.Response, *ParseResult, *ErrorResponse, error) {
	if url == "" {
		return nil, nil, nil, NewArgError("url", "can't be empty")
	}
	for _, f := range fields {
		if f != "next" && f != "text" {
			return nil, nil, nil, NewArgError("fields", "now only support next and text options")
		}
	}
	path := fmt.Sprintf("%s%s", s.client.BaseURL, articleBasePath)

	req, err := http.NewRequest(http.MethodGet, path, nil)

	if err != nil {
		return nil, nil, nil, err
	}

	q := req.URL.Query()
	q.Add("token", s.client.Token)
	q.Add("url", url)

	if fields != nil {
		q.Add("fields", strings.Join(fields, ","))
	}
	req.URL.RawQuery = q.Encode()
	result := new(ParseResult)
	resp, errorResponse, err := s.Do(req, result)
	if err != nil {
		return nil, nil, nil, err
	}

	return resp, result, errorResponse, nil
}

func (s *ArticleServiceOp) Do(req *http.Request, v interface{}) (*http.Response, *ErrorResponse, error) {
	resp, err := s.client.Do(req)

	if err != nil {
		return nil, nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	errorResponse, err := CheckResponse(resp)

	if err != nil {
		return nil, nil, err
	}

	if errorResponse != nil {
		return resp, errorResponse, nil
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	return resp, nil, nil
}

func CheckResponse(r *http.Response) (*ErrorResponse, error) {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil, nil
	}

	errorResponse := &ErrorResponse{}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			return nil, err
		}
	}

	return errorResponse, nil
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	// Response *http.Response
	// Message error message
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

// func (r *ErrorResponse) Error() string {
// 	return fmt.Sprintf("%v %v: %d %v",
// 		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
// }
