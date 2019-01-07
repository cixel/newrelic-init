package main

import "net/http"

func named_arg(a int) {
}

func unnamed_arg(int) {
}

func named_return() (a int) {
	a = 1
	return
}

func unnamed_return() int {
	return 1
}

func main() {
	named_arg(0)
	unnamed_arg(0)
	_ = named_return()
	_ = unnamed_return()

	http.HandleFunc("asdf", nil)

	var r_named_arg = named_arg
	r_named_arg(0)

	var r_httpHandleFunc = http.HandleFunc
	r_httpHandleFunc("asdf", nil)
}
