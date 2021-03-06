[![Build Status](https://travis-ci.com/alastairruhm/go-url2io.svg?branch=master)](https://travis-ci.com/alastairruhm/go-url2io)

# go-url2io

Golang API client for http://url2io.applinzi.com/

## usage example

```go
package main

import (
	"log"
	"strconv"
	"strings"

	url2io "github.com/alastairruhm/go-url2io"
)

func main() {
	client := url2io.NewClient("token-of-url2io", nil)

	// article api

	_, result, err := client.Article.Parse("https://colobu.com/2018/08/27/learn-go-module/", nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("result:\n")
	log.Printf("title: %s\n", result.Title)
	log.Printf("content: %s\n", strconv.Quote(getAbstract(result.Content)))
	log.Printf("url: %s", result.URL)

	// aritcle api with text parameter

	_, result, err = client.Article.Parse("https://colobu.com/2018/08/27/learn-go-module/", []string{"text"})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("result:\n")
	log.Printf("title: %s\n", result.Title)
	log.Printf("text: %s\n", strconv.Quote(getAbstract(result.Text)))
	log.Printf("url: %s", result.URL)
}

func getAbstract(content string) string {
	content = strings.TrimSpace(content)
	if len(content) < 100 {
		return content
	}
	return content[:100] + "..."
}

```

