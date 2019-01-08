package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cixel/newrelic-init/testutil"
)

func TestInjectInit(t *testing.T) {
	tests := map[string]bool{
		"main":    true,
		"nonmain": true,
	}

	for test, expect := range tests {
		t.Run(test, func(t *testing.T) {
			pkg := testutil.LoadPkg(t)

			const appname = "test_app_name"
			const license = "test_license"

			injectInit(pkg)
			buf := fileToBuf(pkg.Syntax[0], "foo", ".")

			str := buf.String()

			if expect != strings.Contains(str, "conf :=") {
				t.Fatalf("missing assignment to conf:\n%s", str)
			}

			if expect != strings.Contains(str, "github.com/newrelic/go-agent") {
				t.Fatalf("missing newrelic import:\n%s", str)
			}

			if expect != strings.Contains(str, fmt.Sprintf(`Getenv("%s")`, nrLicenseEnv)) {
				t.Fatalf("missing expected check for license key env:\n%s", str)
			}

			if expect != strings.Contains(str, fmt.Sprintf(`Getenv("%s")`, nrAppEnv)) {
				t.Fatalf("missing expected check for app name env:\n%s", str)
			}

			if expect != strings.Contains(str, `"os"`) {
				t.Fatalf("missing os import\n %s", str)
			}

			testutil.CompareGolden(t, buf.Bytes())
		})
	}
}
