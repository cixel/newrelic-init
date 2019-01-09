newrelic-init adds New Relic to a package:

```go
var newrelicApp newrelic.Application
func init() {
	conf := newrelic.NewConfig(os.Getenv("NEW_RELIC_APP_NAME"), os.Getenv("NEW_RELIC_LICENSE_KEY"))
	app, _ := newrelic.NewApplication(config)
	newrelicApp = app
}
```

It also attempts to wrap arguments to http.HandleFunc.
See [the New Relic docs](https://docs.newrelic.com/docs/agents/go-agent/installation/install-new-relic-go) for more details.

## Usage

```
$ go get github.com/cixel/newrelic-init
$ newrelic-init [OPTION]... SOURCE
```

Ensure the `NEW_RELIC_APP_NAME` and `NEW_RELIC_LICENSE_KEY` variables are set before running the instrumented app.

`SOURCE` must be valid as a package path.

## TODO

* config for adding transactions to named package functions

## FIXME
* avoid adding init if it's already there
* avoid re-wrapping previously wrapped `http.HandleFunc`
* handle code like `http.HandleFunc(x())` where multi-returns need to be saved in order to be passed into `newrelic.WrapHandleFunc`
