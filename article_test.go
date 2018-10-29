package url2io

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestArticle_ParseContent(t *testing.T) {
	setup()

	defer teardown()

	mux.HandleFunc("/article", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"url": "https://httpbin.org/", "date": "2018-05-18 16:57:48", "content": "test content", "title": "test title"}`)
		testMethod(t, r, http.MethodGet)
	})

	_, result, err := client.Article.Parse("https://httpbin.org/", nil)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if result.Content != "test content" {
		t.Fatalf("unexpected response")
	}
}

func TestArticle_ParseParamURLError(t *testing.T) {
	_, _, err := client.Article.Parse("", nil)

	expected := &ArgError{arg: "url", reason: "can't be empty"}

	if !reflect.DeepEqual(err, expected) {
		t.Errorf("Error = %#v, expected %#v", err, expected)
	}

	if err.Error() != "url is invalid because can't be empty" {
		t.Errorf("unexpected ArgError.Error()")
	}
}
