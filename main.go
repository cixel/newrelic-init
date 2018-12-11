package main

import (
	"flag"
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
