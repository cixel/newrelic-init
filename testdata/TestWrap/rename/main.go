package main

import "net/http"

// var h = http.HandleFunc
var h func(pattern string, handler func(http.ResponseWriter, *http.Request))

type aliasHandlerFunc http.HandlerFunc

func main() {
	h = http.HandleFunc
	var aHandler aliasHandlerFunc = func(http.ResponseWriter, *http.Request) {}

	h("/a", aHandler)
}

// ------

// package main

// import (
// 	"net/http"
// )

// // var h = http.HandleFunc

// type aliasHandlerFunc http.HandlerFunc

// func main() {
// 	var h = http.HandleFunc
// 	var aHandler aliasHandlerFunc = func(http.ResponseWriter, *http.Request) {}

// 	h("/a", aHandler)
// }
