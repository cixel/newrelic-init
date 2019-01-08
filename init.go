package main

import (
	"fmt"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

const newrelicPkgPath = "github.com/newrelic/go-agent"

const nrLicenseEnv = "NEW_RELIC_LICENSE_KEY"
const nrAppEnv = "NEW_RELIC_APP_NAME"

var packageNameHints = map[string]string{newrelicPkgPath: "newrelic"}

func injectInit(pkg *decorator.Package) {
	f := pkg.Syntax[0]

	d := buildInitFunc()

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

// create a call to os.Getenv("env")
func makeGetenv(env string) *dst.CallExpr {
	e := &dst.BasicLit{
		Kind:  token.STRING,
		Value: fmt.Sprintf(`"%s"`, env),
	}

	call := &dst.CallExpr{
		// Instead of a selector expr, let dst handle import resolution for us
		// by using decorated Ident (https://github.com/dave/dst#imports)
		Fun: &dst.Ident{
			Name: "Getenv",
			Path: "os",
		},
		Args: []dst.Expr{e},
	}

	return call
}

func buildInitFunc() dst.Decl {
	// --- line 1 ---
	// config := newrelic.NewConfig(os.Getenv(...), os.Getenv(...))
	callNewConf := &dst.CallExpr{
		Fun: &dst.Ident{Path: newrelicPkgPath, Name: "NewConfig"},
		Args: []dst.Expr{
			makeGetenv(nrAppEnv),
			makeGetenv(nrLicenseEnv),
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
	// app, err := newrelic.NewApplication(config)
	callNewApp := &dst.CallExpr{
		Fun:  &dst.Ident{Path: newrelicPkgPath, Name: "NewApplication"},
		Args: []dst.Expr{dst.NewIdent("conf")},
	}

	assignLocal := &dst.AssignStmt{
		// FIXME: we shouldn't be underscoring the error
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
