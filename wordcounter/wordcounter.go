package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/net/html"
)

// TextExtractor is an interface with an "extract" function that processes
// and refines text from some source.
type TextExtractor interface {
	extract() string
}

// BasicTextExtractor is an extension of TextExtractor for raw strings.
type BasicTextExtractor struct {
	text string
}

// URLTextExtractor is an extension of TextExtractor for a URL.
type URLTextExtractor struct {
	url string
}

// FileTextExtractor is an extension of TextExtractor for a file path.
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
	const DefaultTimeout = 30 * time.Second
	httpClient := &http.Client{Timeout: DefaultTimeout}
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

// GetWordCount returns a map where each key is a word found in the
// input TextExtractor's processed/refined text, and each value is
// the key's corresponding frequency in said text.
func GetWordCount(ex TextExtractor) map[string]int {
	wordCountMap := make(map[string]int)

	words := regexp.MustCompile(`\s+`).Split(ex.extract(), -1)
	for _, word := range words {
		if len(word) == 0 { // Skip whitespaces.
			continue
		}

		// If word is new, then add (word, 1).
		// Otherwise, increment corresponding count by 1.
		freq, prevSeenWord := wordCountMap[word]
		switch prevSeenWord {
		case true:
			wordCountMap[word] = freq + 1
		default:
			wordCountMap[word] = 1
		}
	}

	return wordCountMap
}

// GetAggregateWordCount applies GetWordCount to each TextExtractor in a given
// slice and aggregates all of the results into a single word-count map.
func GetAggregateWordCount(extractors []TextExtractor) map[string]int {
	resultMap := make(map[string]int)

	// Channel for storing extractors' results on-the-go.
	type WordCountPair struct {
		word  string
		count int
	}
	wcpch := make(chan WordCountPair)
	// Channel to keep track of which extractors are done.
	numExtractors := len(extractors)
	done := make(chan int, numExtractors)

	// Initialize a goroutine/thread for each extractor; this thread
	// shall write out each of the extractor's (word, count) pairs to wpch.
	// Once it is done, it shall add an entry to the buffered channel, done.
	for i, ex := range extractors {
		go func() {
			wcm := GetWordCount(ex)
			for word, count := range wcm {
				wcpch <- WordCountPair{word, count}
			}
			done <- i
		}()
	}

	// Monitor done. Once the number of values read from done equals the
	// total number of extractors, we are finished and can close all channels.
	go func() {
		for i := 0; i < numExtractors; i++ {
			<-done
		}
		close(done)
		close(wcpch)
	}()

	// Read the latest (word, count) pair from the channel, and create/update
	// the corresponding entry in resultMap.
	for wcp := range wcpch {
		currentCount, present := resultMap[wcp.word]
		switch present {
		// If word already exists, add received count to current count.
		case true:
			resultMap[wcp.word] = currentCount + wcp.count
		// Otherwise, create an entry for word with received count.
		default:
			resultMap[wcp.word] = wcp.count
		}
	}
	return resultMap
}

func main() {
	// Parse command-line arguments.
	s := flag.String("s", "", "Generate word histogram for given string.")
	u := flag.String("u", "", "Generate word histogram for given URL.")
	f := flag.String("f", "", "Generate word histogram for given file path.")
	flag.Parse()

	// Construct TextExtractors from command-line arguments.
	var extractors []TextExtractor
	if len(*s) > 0 {
		extractors = append(extractors, BasicTextExtractor{*s})
	}
	if len(*u) > 0 {
		urlTokens := regexp.MustCompile(`\s+`).Split(*u, -1)
		for _, url := range urlTokens {
			extractors = append(extractors, URLTextExtractor{url})
		}
	}
	if len(*f) > 0 {
		filePathTokens := regexp.MustCompile(`\s+`).Split(*f, -1)
		for _, filePath := range filePathTokens {
			extractors = append(extractors, FileTextExtractor{filePath})
		}
	}
	// Error out if nothing was specified.
	if len(extractors) == 0 {
		log.Fatal("Please specify at least 1 URL or file path.")
	}

	// Get aggregate word counts for all specified extractors.
	wcm := GetAggregateWordCount(extractors)
	for word, count := range wcm {
		fmt.Printf("%s: %d\n", word, count)
	}
}
