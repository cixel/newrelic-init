package main

import "net/http"

func makeHandler() http.HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {}
}

func main() {
	aHandler := func(http.ResponseWriter, *http.Request) {}

	http.HandleFunc("/a", aHandler)
	http.HandleFunc("/b", func(http.ResponseWriter, *http.Request) {})
	http.HandleFunc("/c", makeHandler())
}
