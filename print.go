package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"log"
)

// NodeToString formats an ast.Node and returns the resulting string
func nodeToString(n ast.Node, fset *token.FileSet) string {
	buf := nodeToBuf(n, fset)
	return buf.String()
}

// formats an ast.Node and returns the resulting buffer
func nodeToBuf(n ast.Node, fset *token.FileSet) bytes.Buffer {
	var buf bytes.Buffer
	err := format.Node(&buf, fset, n)
	if err != nil {
		log.Fatal(err)
	}
	return buf
}
