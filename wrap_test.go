package main

import (
	"testing"

	"github.com/cixel/newrelic-init/testutil"
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
)

func TestWrap(t *testing.T) {
	tests := map[string]bool{
		// "base":   true,
		"rename": true,
		"fake":   false,
	}

	for test := range tests {
		t.Run(test, func(t *testing.T) {
			pkg := testutil.LoadPkg(t)

			for _, file := range pkg.Syntax {
				dstutil.Apply(file, nil, func(c *dstutil.Cursor) bool {
					if call, ok := c.Node().(*dst.CallExpr); ok {
						wrap(pkg, file, call, c)
					}
					return true
				})
			}
		})
	}
}

func TestHandleFunc(t *testing.T) {
}
