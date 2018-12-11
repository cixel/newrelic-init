package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

// TODO: refactor so that this passes in the decl directly rather than
// calling buildInitFunc. could make testing easier.
func injectInit(pkg *packages.Package, name, key string) {
	if pkg.Name != "main" {
		return
	}

	f := pkg.Syntax[0]
	// not sure if I should go through apply/replace here, or just modify the file directly
	astutil.Apply(f, func(c *astutil.Cursor) bool {
		file, ok := c.Node().(*ast.File)
		if !ok {
			return true
		}

		d := buildInitFunc(name, key)

		// add our package-wide 'app' variable to assign to
		// for use by wrappers
		file.Decls = append(file.Decls, appVar())

		// cannot c.Replace File nodes (see Cursor doc for reasoning) so we modify the File directly
		file.Decls = append(file.Decls, d)

		// FIXME
		astutil.AddNamedImport(pkg.Fset, f, "newrelic", "github.com/newrelic/go-agent")

		return false

	}, nil)
}

// creates the node: `var newrelicApp newrelic.Application`
func appVar() *ast.GenDecl {
	d := &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent("newrelicApp")},
				Type: &ast.SelectorExpr{
					X:   ast.NewIdent("newrelic"),
					Sel: ast.NewIdent("Application"),
				},
			},
		},
	}

	return d
}

// config := newrelic.NewConfig("YOUR_APP_NAME", "_YOUR_NEW_RELIC_LICENSE_KEY_")
// app, err := newrelic.NewApplication(config)
// if err != nil {
// }
func buildInitFunc(name, key string) ast.Decl {
	// line 1
	newConf := fmt.Sprintf(`newrelic.NewConfig("%s", "%s")`, name, key)

	callNewConf, err := parser.ParseExpr(newConf)
	if err != nil {
		// 'this should never happen'
		panic(err)
	}

	assignConf := &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent("conf")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			callNewConf,
		},
	}

	// line 2
	newApp := fmt.Sprintf(`newrelic.NewApplication(conf)`)
	callNewApp, err := parser.ParseExpr(newApp)
	if err != nil {
		// FIXME
		panic(err)
	}

	assignApp := &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent("app"), ast.NewIdent("err")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			callNewApp,
		},
	}

	decl := &ast.FuncDecl{
		Name: ast.NewIdent("init"),
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				assignConf,
				assignApp,
			},
		},

		// init functions have no types
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{},
			},
		},
	}

	return decl
}
