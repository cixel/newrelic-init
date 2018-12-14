package testutil

import (
	"path/filepath"
	"testing"

	"github.com/dave/dst/decorator"
	"golang.org/x/tools/go/packages"
)

// LoadPkg loads a package using the tests's name, returning the syntax and
// type info for use in testing. All rewrite tests use data stored in
// testdata/<TestName>/**.go
func LoadPkg(t testing.TB) *decorator.Package {
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  filepath.Join("testdata", t.Name()),
	}

	pkgs, err := decorator.Load(cfg, ".")
	if err != nil {
		t.Fatal(err)
	}

	if len(pkgs) != 1 {
		t.Fatalf("unexpected number of pkgs loaded (have %d, want %d)", len(pkgs), 1)
	}

	return pkgs[0]
}
