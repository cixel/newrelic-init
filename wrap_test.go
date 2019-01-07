package main

import (
	"go/ast"
	"strings"
	"testing"

	"github.com/cixel/newrelic-init/testutil"
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
)

func TestWrap(t *testing.T) {
	tests := map[string]bool{
		"base":   true,
		"rename": true,

		"fake": true,
	}

	for test, expect := range tests {
		t.Run(test, func(t *testing.T) {
			pkg := testutil.LoadPkg(t)

			for _, file := range pkg.Syntax {
				dstutil.Apply(file, nil, func(c *dstutil.Cursor) bool {
					if call, ok := c.Node().(*dst.CallExpr); ok {
						wrap(pkg, call, c)
					}
					return true
				})
			}

			buf := fileToBuf(pkg.Syntax[0], "foo", ".")

			str := buf.String()

			if !expect == strings.Contains(str, "newrelic.WrapHandleFunc") {
				t.Fatalf("missing call to WrapHandleFunc:\n%s", str)
			}

			testutil.CompareGolden(t, buf.Bytes())
		})
	}
}

func TestDoesFuncMatch(t *testing.T) {
	type matchCriteria struct {
		sig            string
		pkgPath        string
		shouldMatch    bool
		shouldMatchPkg bool
	}

	matchtests := map[string]matchCriteria{
		"named_arg": matchCriteria{
			sig:            "func(a int)",
			pkgPath:        "github.com/cixel/newrelic-init/testdata/TestDoesFuncMatch/main",
			shouldMatch:    true,
			shouldMatchPkg: true,
		},
		"unnamed_arg": matchCriteria{
			sig:            "func(int)",
			pkgPath:        "github.com/cixel/newrelic-init/testdata/TestDoesFuncMatch/main",
			shouldMatch:    true,
			shouldMatchPkg: true,
		},
		"named_return": matchCriteria{
			sig:            "func() (a int)",
			pkgPath:        "github.com/cixel/newrelic-init/testdata/TestDoesFuncMatch/main",
			shouldMatch:    true,
			shouldMatchPkg: true,
		},
		"unnamed_return": matchCriteria{
			sig:            "func() int",
			pkgPath:        "github.com/cixel/newrelic-init/testdata/TestDoesFuncMatch/main",
			shouldMatch:    true,
			shouldMatchPkg: true,
		},

		// external
		"http.HandleFunc": matchCriteria{
			sig:            httpHandleFuncType,
			pkgPath:        "net/http",
			shouldMatch:    true,
			shouldMatchPkg: true,
		},

		// renames
		"r_named_arg": matchCriteria{
			sig:            "func(a int)",
			pkgPath:        "github.com/cixel/newrelic-init/testdata/TestDoesFuncMatch/main",
			shouldMatch:    true,
			shouldMatchPkg: true,
		},
		"r_httpHandleFunc": matchCriteria{
			sig:         httpHandleFuncType,
			pkgPath:     "net/http",
			shouldMatch: true,
			// known limitation of matching with pkgPath
			shouldMatchPkg: false,
		},
	}

	t.Run("main", func(t *testing.T) {
		pkg := testutil.LoadPkg(t)

		for _, file := range pkg.Syntax {
			dstutil.Apply(file, nil, func(c *dstutil.Cursor) bool {
				if call, ok := c.Node().(*dst.CallExpr); ok {
					c, ok := pkg.Decorator.Ast.Nodes[call].(*ast.CallExpr)

					if !ok {
						t.Fatal("couldn't map dst.CallExpr to ast.CallExpr")
					}

					fname := astToString(c.Fun)

					if mc, ok := matchtests[fname]; ok {
						m := doesFuncTypeMatch(c, mc.sig, pkg.TypesInfo)
						if mc.shouldMatch != m {
							t.Fatalf("incorrect function match for %s", fname)
						}

						mp := doesFuncTypeMatchWithPkg(c, mc.sig, mc.pkgPath, pkg.TypesInfo)
						if mc.shouldMatchPkg != mp {
							t.Fatalf("incorrect function match for %s (pkg %s)", fname, mc.pkgPath)
						}
					}
				}
				return true
			})
		}
	})
}

func TestHandleFunc(t *testing.T) {
}
