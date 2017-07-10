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

func GetValidDenominations() map[string]string {
	validDenom := make(map[string]string)

	parser := &HTMLParser{Url: "http://www.google.com/finance/converter"}
	parser.HandleToken = func() {
		current := parser.Tokenizer.Token()
		switch {
		case current.Type == html.StartTagToken && current.Data == "option":
			for _, a := range current.Attr {
				if a.Key == "value" {
					parser.Tokenizer.Next() // Read the next token.
					validDenom[a.Val] = parser.Tokenizer.Token().Data
					break
				}
			}
		}
	}
	parser.parseTillEndTag("select")

	return validDenom
}

func ValidateDenominations(inputs ...string) []string {
	// Create []string with all-lowercase, trimmed versions of inputs.
	numInputs := len(inputs)
	processedInputs := make([]string, 0, numInputs)
	for _, input := range inputs {
		processedVal := strings.ToLower(strings.TrimSpace(input))
		if len(processedVal) == 0 {
			log.Fatal("Invalid denomination specified!")
		}
		processedInputs = append(processedInputs, processedVal)
	}

	// For each processed input string, try to find a matching short or
	// long name in the set of valid denominations.
	allValidDenoms := GetValidDenominations()
	specifiedDenoms := make([]string, 0, numInputs)
	for _, in := range processedInputs {
		var target *string = nil
		for shortName, longName := range allValidDenoms {
			sn, ln := strings.ToLower(shortName), strings.ToLower(longName)
			// Check for strict match with short or loose with long.
			if sn == in || strings.Contains(ln, in) {
				target = &shortName
				break
			}
		}
		// Error out if we can't find a match.
		if target == nil {
			log.Fatalf("Invalid denomination: %s!", in)
		}
		// Otherwise, add match's shortName to specifiedDenoms.
		specifiedDenoms = append(specifiedDenoms, *target)
	}
	return specifiedDenoms
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
	listDenom := flag.Bool("l", false, "List all valid denominations.")
	flag.Parse()

	// "-l" takes precedence over all other flags.
	switch (*listDenom) {
	case true:
		for shortName, longName := range GetValidDenominations() {
			fmt.Printf("%s: %s\n", shortName, longName)
		}

	default:
		// Validate parameters.
		if (*amount < 0) {
			log.Fatal("Amount cannot be negative!")
		}
		validated := ValidateDenominations(*from, *to)
		f, t := validated[0], validated[1]
		// Perform conversion, and print result.
		fmt.Printf("%d %s = %s %s\n", *amount, f, Convert(*amount, f, t), t)
	}
}