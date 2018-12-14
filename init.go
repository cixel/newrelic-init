package main

import (

	// "go/dst"

	"fmt"
	"go/token"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

const newrelicPkgPath = "github.com/newrelic/go-agent"

// TODO: refactor so that this passes in the decl directly rather than
// calling buildInitFunc. could make testing easier.
func injectInit(pkg *decorator.Package, name, key string) {
	if pkg.Name != "main" {
		return
	}

	f := pkg.Syntax[0]

	d := buildInitFunc(name, key)

	// We cannot c.Replace File nodes (see Cursor doc for reasoning)
	// so we modify the File directly.
	// First add our package-wide 'app' variable (for use by wrappers) to assign to.
	// Then add the init function.
	f.Decls = append(f.Decls, appVar())
	f.Decls = append(f.Decls, d)

	addImport(f, "newrelic", newrelicPkgPath)
}

func stripQuotes(s string) string {
	return strings.Replace(s, "\"", "", -1)
}

func addImport(f *dst.File, name, path string) {
	var added bool
	for _, imp := range f.Imports {
		if stripQuotes(imp.Path.Value) == path {
			if imp.Name == nil || imp.Name.Name != name {
				imp.Name = dst.NewIdent(name)
				added = true
			}
		}
	}

	if !added {
		imp := &dst.GenDecl{
			Tok: token.IMPORT,
			Specs: []dst.Spec{
				&dst.ImportSpec{
					Name: dst.NewIdent(name),
					Path: &dst.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"%s"`, path),
					},
				},
			},
		}
		f.Decls = append([]dst.Decl{imp}, f.Decls...)
	}
}

// creates the node: `var newrelicApp newrelic.Application`
func appVar() *dst.GenDecl {
	d := &dst.GenDecl{
		Tok: token.VAR,
		Specs: []dst.Spec{
			&dst.ValueSpec{
				Names: []*dst.Ident{dst.NewIdent("newrelicApp")},
				Type: &dst.SelectorExpr{
					X:   dst.NewIdent("newrelic"),
					Sel: dst.NewIdent("Application"),
				},
			},
		},
	}

	return d
}

func buildInitFunc(name, key string) dst.Decl {
	// --- line 1 ---
	// config := newrelic.NewConfig("YOUR_APP_NAME", "_YOUR_NEW_RELIC_LICENSE_KEY_")
	callNewConf := &dst.CallExpr{
		Fun: &dst.SelectorExpr{
			X:   dst.NewIdent("newrelic"),
			Sel: dst.NewIdent("NewConfig"),
		},
		Args: []dst.Expr{
			&dst.BasicLit{
				Kind:  token.STRING,
				Value: `"` + name + `"`,
			},
			&dst.BasicLit{
				Kind:  token.STRING,
				Value: `"` + key + `"`,
			},
		},
	}

	assignConf := &dst.AssignStmt{
		Lhs: []dst.Expr{dst.NewIdent("conf")},
		Tok: token.DEFINE,
		Rhs: []dst.Expr{
			callNewConf,
		},
	}

	// ---line 2---
	// TODO: we shouldn't be underscoring the error
	// app, err := newrelic.NewApplication(config)
	callNewApp := &dst.CallExpr{
		Fun: &dst.SelectorExpr{
			X:   dst.NewIdent("newrelic"),
			Sel: dst.NewIdent("NewApplication"),
		},
		Args: []dst.Expr{dst.NewIdent("conf")},
	}

	assignLocal := &dst.AssignStmt{
		Lhs: []dst.Expr{dst.NewIdent("app"), dst.NewIdent("_")},
		Tok: token.DEFINE,
		Rhs: []dst.Expr{
			callNewApp,
		},
	}

	// ---line 3---
	// newrelicApp = app
	assignPkg := &dst.AssignStmt{
		Lhs: []dst.Expr{dst.NewIdent("newrelicApp")},
		Tok: token.ASSIGN,
		Rhs: []dst.Expr{dst.NewIdent("app")},
	}

	decl := &dst.FuncDecl{
		Name: dst.NewIdent("init"),
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				assignConf,
				assignLocal,
				assignPkg,
			},
		},

		// (empty type)
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{},
			},
			Results: &dst.FieldList{
				List: []*dst.Field{},
			},
		},
	}

	return decl
}
