package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	_, _ = fmt.Fprintln(os.Stderr, `Usage: instrument [OPTION]... SOURCE

newrelic-init adds New Relic to a package:

var newrelicApp newrelic.Application
func init() {
	conf := newrelic.NewConfig(os.Getenv("NEW_RELIC_APP_NAME"), os.Getenv("NEW_RELIC_LICENSE_KEY"))
	app, _ := newrelic.NewApplication(config)
	newrelicApp = app
}

It also attempts to wrap arguments to http.HandleFunc.
See https://docs.newrelic.com/docs/agents/go-agent/installation/install-new-relic-go for more details.

Ensure the NEW_RELIC_APP_NAME and NEW_RELIC_LICENSE_KEY variables are set before running the instrumented app.

SOURCE must be valid as a package path.

Flags:`)
	flag.PrintDefaults()
}
