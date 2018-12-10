package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

func main() {
	var (
		name  string
		key   string
		write bool
	)

	flag.StringVar(&name, "name", "", "newrelic account name")
	flag.StringVar(&key, "key", "", "newrelic license key")
	flag.BoolVar(&write, "w", false, "write to source file instead of stdout")
	flag.Usage = usage
	flag.Parse()
	if len(flag.Args()) != 1 {
		usage()
		// http://tldp.org/LDP/abs/html/exitcodes.html
		os.Exit(2)
	}

	log.SetPrefix("newrelic-init ")
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	src := flag.Args()[0]

	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		Dir:  src,
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		log.Fatal(err)
	}

	packages.Visit(pkgs, func(*packages.Package) bool {
		// we only want to visit the given package,
		// but we want to type check on everything
		return false
	}, func(pkg *packages.Package) {
		for _, f := range pkg.Syntax {
			astutil.Apply(f, nil, func(c *astutil.Cursor) bool {
				return true
			})
		}
	})
}

func injectInit(pkg *packages.Package, name, key string) {
	if pkg.Name != "main" {
		return
	}

	f := pkg.Syntax[0]
	// not sure if I should go through apply/replace here, or just modify the file directly
	astutil.Apply(f, func(c *astutil.Cursor) bool {
		if file, ok := c.Node().(*ast.File); ok {
			d := buildInitFunc(name, key)

			// cannot c.Replace File nodes (see Cursor doc for reasoning) so we modify the File directly
			file.Decls = append(file.Decls, d)

			// FIXME
			astutil.AddImport(pkg.Fset, f, "github.com/newrelic/go-agent")
			return false
		}
		return true
	}, nil)
}

func buildInitFunc(name, key string) ast.Decl {
	newConf := fmt.Sprintf(`newrelic.NewConfig("%s", "%s")`, name, key)

	expr, err := parser.ParseExpr(newConf)
	if err != nil {
		// 'this should never happen'
		panic(err)
	}

	decl := &ast.FuncDecl{
		Name: ast.NewIdent("init"),
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("conf")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						expr,
					},
				},
			},
		},

		// empty
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
