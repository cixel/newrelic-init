package main

import "net/http"

func fake(pattern string, handler func(http.ResponseWriter, *http.Request)) {
}

func main() {
	fake("/a", func(http.ResponseWriter, *http.Request) {})
}
