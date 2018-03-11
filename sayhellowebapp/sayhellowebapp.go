package main

import (
	"fmt"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.URL.Path[1:])
	if len(name) == 0 {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!", name)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}