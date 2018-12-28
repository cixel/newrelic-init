package main

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

const httpHandleFuncType = "func(pattern string, handler func(net/http.ResponseWriter, *net/http.Request))"

func wrap(pkg *decorator.Package, f *dst.File, c *dst.CallExpr, cursor *dstutil.Cursor) {
	call, ok := pkg.Decorator.Ast.Nodes[c].(*ast.CallExpr)

	if !ok {
		return
	}

	file, ok := pkg.Decorator.Ast.Nodes[f].(*ast.File)
	if !ok {
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!")
		return
	}

	// prog, pkgs := ssautil.Packages([]*packages.Package{pkg.Package}, ssa.PrintPackages)
	// prog, pkgs := ssautil.Packages([]*packages.Package{pkg.Package}, ssa.PrintFunctions)
	prog, pkgs := ssautil.Packages([]*packages.Package{pkg.Package}, 0)

	// _ = prog

	// p := pkgs[0]
	// v := p.Var("h")
	// fmt.Println("members:", p.Members)
	// if v != nil {
	// 	fmt.Println("!!!!")
	// 	fmt.Println(v)
	// 	fmt.Printf("obj: %+v\n", v.Object())
	// 	fmt.Printf("refs: %+v\n", v.Referrers())
	// 	fmt.Printf("ops: %+v\n", v.Operands([]*ssa.Value{}))
	// 	// fmt.Printf("ops: %+v\n", v.Operands([]*ssa.Value{}))
	// 	fmt.Printf("pkg: %+v\n", v.Package())
	// }

	doesFuncMatch(call, httpHandleFuncType, "net/http", pkg.TypesInfo, prog, pkgs[0], file) // FIXME check length
	// typ, ok := pkg.TypesInfo.TypeOf(call.Fun).(*types.Signature)
}

// TODO this will:
// if sigs don't match match --> FALSE
// i := extractIdent(call.Fun)
// if i == nil --> FALSE
// get Uses[i]
// get Pkg()
// if Pkg() matches --> TRUE
// else get parent scope
// use Object ID to find corresponding scope member
// re-run on that scope member
//
// idea is just to walk up the scope until we can't anymore and try to
// match the package path
func doesFuncMatch(call *ast.CallExpr, sig, pkgPath string, info *types.Info, prog *ssa.Program, pkg *ssa.Package, file *ast.File) bool {
	typ, ok := info.TypeOf(call.Fun).(*types.Signature)
	if !ok {
		return false
	}

	if typ.String() != sig {
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

	refs, _ := astutil.PathEnclosingInterval(file, id.Pos(), id.End())
	fmt.Println(refs)
	v, ok := use.(*types.Var)
	if !ok {
		fmt.Println("heck")
		return false
	}

	val, addr := prog.VarValue(v, pkg, refs)
	fmt.Printf("val: %+v\n", val)
	fmt.Printf("addr: %+v\n", addr)
	fmt.Printf("pkg: %+v\n", pkg)
	// fmt.Printf("prog: %+v\n", prog)

	enc := ssa.EnclosingFunction(pkg, refs)
	fmt.Printf("enc: %+v\n", enc)
	// fmt.Printf("enc: %#v\n", enc)

	// m := pkg.Func("main")
	// fmt.Printf("main: %+v\n", m)
	// fmt.Printf("blocks: %+v\n", m.Blocks)
	// vfe, a := m.ValueForExpr(id)
	// fmt.Printf("vfe: %+v\n", vfe)
	// fmt.Printf("a: %+v\n", a)
	// // fmt.Printf("prog: %+v\n", prog.)

	// p := prog.Package(info.)
	// fmt.Println(prog.)

	return isUseInPkg(use, pkgPath)
}

func isUseInPkg(use types.Object, pkgPath string) bool {
	// scope := use.Parent()
	if use.Pkg().Path() == pkgPath {
		return true
	}

	fmt.Println("--")

	fmt.Printf("use: %+v\n", use)
	// scope := use.Parent()
	// fmt.Printf("scope: %+v\n", scope)
	// fmt.Printf("scope str: %+v\n", scope.String())
	// fmt.Printf("pos: %+v\n", use.Pos())
	// fmt.Printf("fpos: %+v\n", fset.Position(use.Pos()))
	// fmt.Println(fset)

	return true
	// return isUseInPkg(def, pkgPath)
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
