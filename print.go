package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"log"
)

// NodeToString formats an ast.Node and returns the resulting string
func nodeToString(n ast.Node) string {
	buf := nodeToBuf(n)
	return buf.String()
}

// formats an ast.Node and returns the resulting buffer
func nodeToBuf(n ast.Node) bytes.Buffer {
	var buf bytes.Buffer
	err := format.Node(&buf, token.NewFileSet(), n)
	if err != nil {
		log.Fatal(err)
	}
	return buf
}
