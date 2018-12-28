package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"log"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/gopackages"
)

// NodeToString formats an ast.Node and returns the resulting string
func filetoString(f *dst.File, path, dir string) string {
	buf := fileToBuf(f, path, dir)
	return buf.String()
}

// formats an ast.Node and returns the resulting buffer
func fileToBuf(f *dst.File, path, dir string) bytes.Buffer {
	var buf bytes.Buffer
	d := decorator.NewRestorerWithImports(path, gopackages.WithHints(dir, packageNameHints))
	err := d.Fprint(&buf, f)
	if err != nil {
		log.Fatal(err)
	}
	return buf
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
