package main

import (
	"testing"

	"github.com/cixel/newrelic-init/testutil"
)

func TestInjectInit(t *testing.T) {
	tests := []string{
		"main",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			testutil.CompareGolden(t, []byte("s"))
		})
	}
}
