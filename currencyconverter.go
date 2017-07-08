package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strings"
	"time"
)

type HTMLParser struct {
	Url string
	Tokenizer *html.Tokenizer
	HandleToken func()
	Stop func() bool
}

func (parser *HTMLParser) parse() {
	// Connect to URL.
	const DEFAULT_TIMEOUT = 30 * time.Second
	httpClient := &http.Client{Timeout: DEFAULT_TIMEOUT}
	response, err := httpClient.Get(parser.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Stream and parse URL's HTML content.
	parser.Tokenizer = html.NewTokenizer(response.Body)
	for !parser.Stop() {
		parser.HandleToken()
	}
	parser.Tokenizer = nil // Clear once we're done.
}

func (parser *HTMLParser) parseTillEndTag(tag string) {
	parser.Stop = func() bool {
		nextTokenType := parser.Tokenizer.Next() // Read the next HTML token.
		switch nextTokenType {
		case html.ErrorToken: // EOF
			return true
		case html.EndTagToken:
			return parser.Tokenizer.Token().Data == tag // </tag>
		default:
			return false
		}
	}
	parser.parse()
}

func GetValidDenominations() []string {
	var validDenom []string

	parser := &HTMLParser{Url: "http://www.google.com/finance/converter"}
	parser.HandleToken = func() {
		current := parser.Tokenizer.Token()
		switch {
		case current.Type == html.StartTagToken && current.Data == "option":
			for _, a := range current.Attr {
				if a.Key == "value" {
					validDenom = append(validDenom, a.Val)
					break
				}
			}
		}
	}
	parser.parseTillEndTag("select")

	return validDenom
}

func AreDenominationsValid(possibleDenom ...string) bool {
	validDenom := strings.Join(GetValidDenominations(), "")
	for _, denom := range possibleDenom {
		if !strings.Contains(validDenom, denom) {
			return false
		}
	}
	return true
}

func Convert(amt int, from, to string) string {
	const URL = "http://www.google.com/finance/converter"
	url := fmt.Sprintf("%s?a=%d&from=%s&to=%s", URL, amt, from, to)
	parser := &HTMLParser{Url: url}

	var result string
	parser.HandleToken = func() {
		current := parser.Tokenizer.Token()
		switch {
		case current.Type == html.StartTagToken && current.Data == "span":
			for _, a := range current.Attr {
				if a.Key == "class" && a.Val == "bld" {
					parser.Tokenizer.Next() // Read the next token.
					result = parser.Tokenizer.Token().Data
					break
				}
			}
		}
	}
	parser.parseTillEndTag("span")
	return strings.TrimSpace(strings.Replace(result, to, "", -1))
} 

func main() {
	// Process command-line arguments.
	amount := flag.Int("a", 1, "Amount being converted.")
	from := flag.String("f", "", "Denomination we are converting from.")
	to := flag.String("t", "", "Denomination we are converting to.")
	flag.Parse()

	// Remove leading/trailing whitespace from "from" and "to."
	f, t := strings.TrimSpace(*from), strings.TrimSpace(*to)

	// Validate parameters.
	if (*amount < 0) {
		log.Fatal("Amount cannot be negative!")
	}
	if !AreDenominationsValid(f, t) {
		log.Fatal("From/To denominations are invalid!")
	}

	// Perform conversion, and print result.
	fmt.Printf("%d %s = %s %s\n", *amount, f, Convert(*amount, f, t), t)
}