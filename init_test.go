package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cixel/newrelic-init/testutil"
	"github.com/dave/dst/decorator"
)

func TestInjectInit(t *testing.T) {
	tests := map[string]bool{
		"main":    true,
		"nonmain": false,
	}

	for test, expect := range tests {
		t.Run(test, func(t *testing.T) {
			pkg := testutil.LoadPkg(t)

			const appname = "test_app_name"
			const license = "test_license"

			injectInit(pkg, appname, license)
			buf := fileToBuf(pkg.Syntax[0])

			testutil.CompareGolden(t, buf.Bytes())

			str := buf.String()

			if expect == !strings.Contains(str, "conf :=") {
				t.Fatalf("missing assignment to conf:\n%s", str)
			}

			if expect == !strings.Contains(str, appname) {
				t.Fatalf("missing app name:\n%s", str)
			}

			if expect == !strings.Contains(str, appname) {
				t.Fatalf("missing license:\n%s", str)
			}

			if expect == !strings.Contains(str, "github.com/newrelic/go-agent") {
				t.Fatalf("missing newrelic import:\n%s", str)
			}
		})
	}
}

func TestAddImport(t *testing.T) {
	tests := map[string]string{
		"not imported":           "package main",
		"imported and named":     fmt.Sprintf("package main\nimport blah \"%s\"", newrelicPkgPath),
		"imported and not named": fmt.Sprintf("package main\nimport \"%s\"", newrelicPkgPath),
		"not imported with docs": "//doc1\n\n//doc2\npackage main",
	}

	for test, code := range tests {
		t.Run(test, func(*testing.T) {
			file, err := decorator.Parse(code)
			if err != nil {
				t.Fatal(err)
			}

			addImport(file, "newrelic", newrelicPkgPath)
			i := fmt.Sprintf(`import newrelic "%s"`, newrelicPkgPath)

			str := filetoString(file)
			c := strings.Count(str, i)
			if c != 1 {
				t.Fatalf("instances of %s in file == %d\n%s", i, c, str)
			}
		})
	}
}
