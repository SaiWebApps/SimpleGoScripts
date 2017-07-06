package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

type TextExtractor interface {
	extract() string
}
type BasicTextExtractor struct {
	text string
}
type URLTextExtractor struct {
	url string
}
type FileTextExtractor struct {
	path string
}

func removeNonWordChars(rawText string) string {
	r := regexp.MustCompile(`\W`)
	replaceWithSpace := func(string) string {
		return " "
	}
	return r.ReplaceAllStringFunc(rawText, replaceWithSpace)
}

func (ex BasicTextExtractor) extract() string {
	return removeNonWordChars(ex.text)
}

func (ex URLTextExtractor) extract() string {
	// Connect to URL.
	const DEFAULT_TIMEOUT = 30 * time.Second
	httpClient := &http.Client{Timeout: DEFAULT_TIMEOUT}
	response, err := httpClient.Get(ex.url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Load web page's HTML into a parser, and return only text b/w tags.
	var buffer bytes.Buffer
	parser := html.NewTokenizer(response.Body)
	for t := parser.Next(); t != html.ErrorToken; t = parser.Next() {
		switch {
		case t == html.TextToken:
			buffer.WriteString(parser.Token().Data)
		}
	}
	return removeNonWordChars(buffer.String())
}

func (ex FileTextExtractor) extract() string {
	bytes, err := ioutil.ReadFile(ex.path)
	if err != nil {
		log.Fatal(err)
	}
	return removeNonWordChars(string(bytes))
}

func GetWordCount(ex TextExtractor) map[string]int {
	wordCountMap := make(map[string]int)

	words := regexp.MustCompile(`\s+`).Split(ex.extract(), -1)
	for _, word := range words {
		// If word is new, then add (word, 1).
		// Otherwise, increment corresponding count by 1.
		freq, prevSeenWord := wordCountMap[word]
		switch (prevSeenWord) {
		case true:
			wordCountMap[word] = freq + 1
		default:
			wordCountMap[word] = 1
		}
	}

	return wordCountMap
}

func main() {
	const URL_ARG_HELP = "Generate word histogram for given URL."
	const DEFAULT_URL = "https://sherlock-holm.es/stories/html/cnus.html"
	url := flag.String("u", DEFAULT_URL, URL_ARG_HELP)
	flag.Parse()

	wcm := GetWordCount(URLTextExtractor{*url})
	for word, count := range wcm {
		fmt.Printf("%s: %d\n", word, count)
	}
}