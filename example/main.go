package main

import (
	"log"
	"strconv"

	url2io "github.com/alastairruhm/go-url2io"
)

func main() {
	client := url2io.NewClient("zF1lMUyCRVuuPaNSqb68yQ", nil)

	// article api
	_, result, errorResponse, err := client.Article.Parse("http://url2io.applinzi.com/docs#url2article", nil)

	if err != nil {
		log.Fatal(err)
	}

	if errorResponse != nil {
		log.Printf("request error\n")
		log.Printf("status code: %d\n", errorResponse.Code)
		log.Printf("error type: %s\n", errorResponse.ErrorType)
		log.Printf("message: %s\n", errorResponse.Message)
		log.Printf("url: %s\n", errorResponse.URL)
		return
	}

	log.Printf("result:\n")
	log.Printf("title: %s\n", result.Title)
	log.Printf("content: %s\n", strconv.Quote(result.Content[:20]))
	log.Printf("url: %s", result.URL)

	// aritcle api with text parameter

	// _, result, errorResponse, err := client.Article.Parse("http://url2io.applinzi.com/docs#url2article", nil)

}
