package main

import (
	"strings"
	"testing"

	"github.com/cixel/newrelic-init/testutil"
)

func TestInjectInit(t *testing.T) {
	tests := []string{
		"main",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			pkg := testutil.LoadPkg(t)

			const appname = "test_app_name"
			const license = "test_license"

			injectInit(pkg, appname, license)
			buf := nodeToBuf(pkg.Syntax[0], pkg.Fset)

			testutil.CompareGolden(t, buf.Bytes())

			str := buf.String()

			if !strings.Contains(str, "conf :=") {
				t.Fatalf("missing assignment to conf:\n%s", str)
			}

			if !strings.Contains(str, appname) {
				t.Fatalf("missing app name:\n%s", str)
			}

			if !strings.Contains(str, appname) {
				t.Fatalf("missing license:\n%s", str)
			}

			if !strings.Contains(str, "github.com/newrelic/go-agent") {
				t.Fatalf("missing newrelic import:\n%s", str)
			}
		})
	}
}
