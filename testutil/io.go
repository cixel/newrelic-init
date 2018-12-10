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

// ReadGolden reads the golden file for the given test
func ReadGolden(t testing.TB) (string, func(string)) {

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(cwd, "testdata", t.Name()+".golden")

	return ReadTestData(t, path), func(s string) {
		if !*update {
			return
		}

		t.Log("updating golden file at", path)
		err := ioutil.WriteFile(path, []byte(s), 0666)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// CompareGolden reads the golden file for the given test
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

	if *update {
		t.Log("updating golden file at", path)
		_, err = file.Write(d)
		if err != nil {
			t.Fatal(err)
		}
	}

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(contents, d) {
		t.Fatalf("expected:\n%s, got:\n%s", contents, d)
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
