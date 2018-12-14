package main

import (
	"bytes"
	"log"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

// NodeToString formats an ast.Node and returns the resulting string
func filetoString(f *dst.File) string {
	buf := fileToBuf(f)
	return buf.String()
}

// formats an ast.Node and returns the resulting buffer
func fileToBuf(f *dst.File) bytes.Buffer {
	var buf bytes.Buffer
	err := decorator.Fprint(&buf, f)
	if err != nil {
		log.Fatal(err)
	}
	return buf
}
