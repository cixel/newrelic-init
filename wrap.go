package main

import (
	"go/ast"
	"go/types"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
)

const httpHandleFuncType = "func(pattern string, handler func(net/http.ResponseWriter, *net/http.Request))"

func wrap(pkg *decorator.Package, c *dst.CallExpr, cursor *dstutil.Cursor) {
	call, ok := pkg.Decorator.Ast.Nodes[c].(*ast.CallExpr)

	if !ok {
		return
	}

	if doesFuncTypeMatch(call, httpHandleFuncType, pkg.TypesInfo) {
		wrapped := &dst.CallExpr{
			Fun: c.Fun,
			Args: []dst.Expr{
				&dst.CallExpr{
					Fun: &dst.SelectorExpr{
						X:   dst.NewIdent("newrelic"),
						Sel: dst.NewIdent("WrapHandleFunc"),
					},
					Args: []dst.Expr{
						dst.NewIdent("newrelicApp"),
						c.Args[0],
						c.Args[1],
					},
				},
			},
		}
		cursor.Replace(wrapped)
	}
}

func doesFuncTypeMatchWithPkg(call *ast.CallExpr, sig, pkgPath string, info *types.Info) bool {
	if !doesFuncTypeMatch(call, sig, info) {
		return false
	}

	id := extractIdent(call.Fun)
	if id == nil {
		return false
	}

	use, ok := info.Uses[id]
	if !ok {
		return false
	}

	if use.Pkg().Path() == pkgPath {
		return true
	}

	return false
}

func doesFuncTypeMatch(call *ast.CallExpr, sig string, info *types.Info) bool {
	typ, ok := info.TypeOf(call.Fun).(*types.Signature)
	if !ok {
		return false
	}

	if typ.String() != sig {
		return false
	}

	return true
}

func extractIdent(expr ast.Expr) *ast.Ident {
	// TODO probably more weird cases
	switch e := expr.(type) {
	// case *ast.ParenExpr:
	case *ast.SelectorExpr:
		return e.Sel
	case *ast.Ident:
		return e
	}

	return nil
}

func wrapHandleFunc() {
}
