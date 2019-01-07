package main

import (
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

const newrelicPkgPath = "github.com/newrelic/go-agent"

var packageNameHints = map[string]string{newrelicPkgPath: "newrelic"}

func injectInit(pkg *decorator.Package, name, key string) {
	f := pkg.Syntax[0]

	d := buildInitFunc(name, key)

	// We cannot c.Replace File nodes (see Cursor doc for reasoning)
	// so we modify the File directly.
	// First add our package-wide 'app' variable (for use by wrappers) to assign to.
	// Then add the init function.
	f.Decls = append(f.Decls, appVar())
	f.Decls = append(f.Decls, d)
}

// creates the node: `var newrelicApp newrelic.Application`
func appVar() *dst.GenDecl {
	d := &dst.GenDecl{
		Tok: token.VAR,
		Specs: []dst.Spec{
			&dst.ValueSpec{
				Names: []*dst.Ident{dst.NewIdent("newrelicApp")},
				Type:  &dst.Ident{Path: newrelicPkgPath, Name: "Application"},
			},
		},
	}

	return d
}

func buildInitFunc(name, key string) dst.Decl {
	// --- line 1 ---
	// config := newrelic.NewConfig("YOUR_APP_NAME", "_YOUR_NEW_RELIC_LICENSE_KEY_")
	callNewConf := &dst.CallExpr{
		Fun: &dst.Ident{Path: newrelicPkgPath, Name: "NewConfig"},
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
		Fun:  &dst.Ident{Path: newrelicPkgPath, Name: "NewApplication"},
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
