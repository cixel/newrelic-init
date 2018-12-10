package testutil

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

// CompareGolden reads a golden file for a given test and compares
// to the given bytes
func CompareGolden(t testing.TB, d []byte) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(cwd, "testdata", t.Name()+".golden")

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = file.Close()
	}()

	var golden []byte

	if *update {
		t.Log("updating golden file at", path)
		_, err = file.Write(d)
		golden = d
		if err != nil {
			t.Fatal(err)
		}
	} else {
		golden, err = ioutil.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}
	}

	if !bytes.Equal(golden, d) {
		c := bytes.Compare(golden, d)
		t.Fatalf("difference at byte %d\nwant:\n%s\ngot:\n%s", c, golden, d)
	}
}

// ReadTestData reads data for a fixture given a set of paths to join
func ReadTestData(t testing.TB, paths ...string) string {
	loc := filepath.Join(paths...)

	f, err := ioutil.ReadFile(loc)
	if err != nil {
		t.Fatal(err)
	}

	return string(f)
}
