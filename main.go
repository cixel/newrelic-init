package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/gopackages"
	"github.com/dave/dst/dstutil"
	"golang.org/x/tools/go/packages"
)

type config struct {
	name  string
	key   string
	write bool
}

func parseFlags() config {
	flag.Usage = usage

	conf := config{}

	flag.StringVar(&conf.name, "name", "", "newrelic account name")
	flag.StringVar(&conf.key, "key", "", "newrelic license key")
	flag.BoolVar(&conf.write, "w", false, "write to source file instead of stdout")

	flag.Parse()

	if len(flag.Args()) != 1 {
		usage()
		// http://tldp.org/LDP/abs/html/exitcodes.html
		os.Exit(2)
	}

	return conf
}

func newrelic(conf config, dir string) *decorator.Package {
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		Dir:  dir,
	}

	pkgs, err := decorator.Load(cfg, ".")
	if err != nil {
		log.Fatal(err)
	}

	pkg := pkgs[0]

	for _, file := range pkg.Syntax {
		dstutil.Apply(file, nil, func(c *dstutil.Cursor) bool {
			if call, ok := c.Node().(*dst.CallExpr); ok {
				wrap(pkg, call, c)
			}
			return true
		})
	}

	injectInit(pkg, conf.name, conf.key)

	return pkg
}

func main() {
	conf := parseFlags()

	log.SetPrefix("newrelic-init ")
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	dir := flag.Args()[0]

	pkg := newrelic(conf, dir)

	if conf.write {
		r := gopackages.WithHints(pkg.Dir, packageNameHints)
		err := pkg.SaveWithResolver(r)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		for _, file := range pkg.Syntax {
			s := filetoString(file, pkg.PkgPath, dir)
			fmt.Println(s)
		}
	}
}
