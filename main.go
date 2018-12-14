package main

import (
	"flag"
	"log"
	"os"

	"github.com/dave/dst/decorator"
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

	pkgs, err := decorator.Load(cfg, ".")
	if err != nil {
		log.Fatal(err)
	}

	pkg := pkgs[0]
	injectInit(pkg, name, key)

	if write {
		err := pkg.Save()
		if err != nil {
			log.Fatal(err)
		}
	}
}
