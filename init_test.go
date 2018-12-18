package main

import (
	"strings"
	"testing"

	"github.com/cixel/newrelic-init/testutil"
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
			buf := fileToBuf(pkg.Syntax[0], "foo", ".")

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
