# Plugin Guide

> a developer's guide to Kanali plugins

## Overview

Want to create you own plugin for *Kanali*? No problem! With *Kanali's* extensible plugin system, it's easy! Here are some steps to follow:

* [Step 1: Create](#step-1-create)
* [Step 2: Test](#step-2-test)
* [Step 3: Compile](#step-3-compile)
* [Step 4: Place](#step-4-place)
* [Step 5: Use](#step-5-use)
* [Step 6 (optional): Version](#step-6-optional-version)

## Step 1: Create

There are essentially only two main design requirements for a Kanali plugin:

1. The plugin must implement Kanali's [`Plugin`](https://github.com/northwesternmutual/kanali/blob/master/plugins/plugin.go) interface.
2. The name of the exported variable implementing the interface must be `Plugin`

To make this simple, I've provided a [Yeoman](http://yeoman.io/) template that will dynamically scaffold everything out for you. You can find the template [here](https://github.com/northwesternmutual/kanali-plugin-template) or just follow the instructions below:

```sh
$ npm install -g yo
$ npm install -g generator-kanali-plugin
$ yo kanali-plugin
```

Every plugin has the ability to intercept a request at two points during the request lifecycle. It does this by invoking the two lifecycle methods that make up the [`Plugin`](https://github.com/northwesternmutual/kanali/blob/master/plugins/plugin.go) interface:

1. `OnRequest` is invoked *before* the request is proxied upstream.
2. `OnResponse` is invoked *after* the request has been returned from the upstream service but *before* the request is returned to the client.

The return value for both of these methods depends on whether an error was encountered or not. If an error was encountered during the plugin logic that should result in the termination of the request, return an [`error`](https://golang.org/pkg/errors/#pkg-examples). If no error was encountered and you would like the request's lifecycle to proceed as normal, return `nil`.

In Go, an [`error`](https://golang.org/pkg/errors/#pkg-examples) is simply an interface. Hence, if you would like to specify an HTTP status code corresponding to your specific error message, the following `type`, which implements this interface, is provided for you to use. If this type is not used, `http.StatusInternalServerError` will be used. An example showing how to use the following type is provided in the template.

```go
type StatusError struct {
	Code int
	Err  error
}
```

A plugin has the ability to define any configuration items it might require. This allows a plugin to be configured via a cli flag, environment variable, or a configuration file. Reference the configuration [documentation](https://github.com/northwesternmutual/kanali#usage-and-configuration) for details. Below is an example of how the [API key plugin](https://github.com/northwesternmutual/kanali-plugin-apikey) uses this feature:

```go
var flagPluginsAPIKeyHeaderKey = config.Flag{
	Long:  "plugins.apiKey.header_key", // please use the convention of plugins.<plugin name>.<configuration name>
	Short: "",
	Value: "apikey",
	Usage: "Name of the HTTP header holding the apikey.",
}

func init() {
	config.Flags.Add(flagPluginsAPIKeyHeaderKey)
}
```

A series of method parameters are provided for optional usage inside of your plugin logic. Here is a table providing detailing their purpose:

name        | type             | methods              |mutability |description
------------|------------------|----------------------|-----------|------------
`ctx`        | [`context.Context`](https://golang.org/pkg/context/) | `OnRequest` `OnResponse` | Mutable | Request context.
`m`        | [`*metrics.Metrics`](https://github.com/northwesternmutual/kanali/blob/master/metrics/metrics.go) | `OnRequest` `OnResponse` | Mutable | Holds various requests metrics for analytics.
`proxy`      | [`spec.ApiProxy`](https://github.com/northwesternmutual/kanali/blob/master/spec/apiproxy.go#L20) | `OnRequest` `OnResponse` | Immutable | This parameter gives you access to the `ApiProxy` struct that matched the incoming request.
`ctlr`       | [`controller.Controller`](https://github.com/northwesternmutual/kanali/blob/master/controller/controller.go#L17) | `OnRequest` `OnResponse` | Immutable | This parameter provides a client by which the Kubernetes api may be accessed.
`req`        | [`http.Request`](https://golang.org/pkg/net/http/#Request) | `OnRequest` `OnResponse` | Mutable | This parameter gives you access to the original HTTP request struct.
`resp`       | [`*http.Response`](https://golang.org/pkg/net/http/#Response) | `OnResponse` | Mutable | This parameter will point to the response that was returned from the upstream service. Note that it is mutable allowing for potential changes in a plugin's logic.
`span`       | [`opentracing-go.Span`](https://godoc.org/github.com/opentracing/opentracing-go#Span) | `OnRequest` `OnResponse` | Immutable | This parameter gives you access to the parent tracing span allowing you to add details (tags) to that span and optionally create new spans in the context of this parent span.

## Step 2: Test

No code is complete without ample test coverage! If you are using the [template](https://github.com/northwesternmutual/kanali-plugin-template) to help bootstrap your plugin, your testing framework is already scaffolded for you. Simply run the following commands:

```sh
$ make test
$ make cover
```

## Step 3: Compile

Go plugins are not compiled into Kanali's binary but instead are loaded at runtime. Go expects are plugin to be an ELF shared object file. Here's an example of how to compile your plugin using the `go` cli:

```sh
$ go build -buildmode=plugin -o myCustomPlugin.so myCustomPlugin.go
```

An important note is that plugins can *only* be compiled on linux. Here's an example of how Kanali compiles the apikey plugin in its `Dockerfile`:

```sh
# download apikey plugin
RUN curl -O https://github.com/northwesternmutual/kanali-plugin-apikey/raw/master/apikey.go
# compile plugin
RUN go build -buildmode=plugin -o apiKey.so plugin.go
# Build project
RUN make build
```

## Step 4: Place

During runtime, Kanali will look for compliled plugins in a location specified by the `plugins-location` cli flag with a default location of `/`. Note that for security reasons, relative paths cannot be used.

Depending on your use case, you may want to dynamically load new plugins without having to rebuild Kanali's Docker container. This is easily accomplished by mapping a volume from the host and placing new compiled plugins there.

## Step 5: Use

Adding a custom plugin to an `ApiProxy` is simple:

```yaml
apiVersion: kanali.io/v1
kind: ApiProxy
metadata:
  name: plugin-example
  namespace: application
spec:
  path: /api/v1/plugin-example
  service:
    port: 8080
    name: my-service
  plugins:
  - name: myCustomPlugin
```

*NOTE:* The plugin name _**must**_ match the file name *(case sensitive)* of the compiled plugin from the previous step.

## Step 6 (optional): Version

Kanali gives you the option to version control your plugins. Below is an example of how to use a certain version of your plugin:

```yaml
apiVersion: kanali.io/v1
kind: ApiProxy
metadata:
  name: plugin-example
  namespace: application
spec:
  path: /api/v1/plugin-example
  service:
    port: 8080
    name: my-service
  plugins:
  - name: myCustomPlugin
    version: 1.0.0
```

Here are some naming rules that must be followed to accomplish this:
* The file name for a compiled plugin representing a specific version must be the plugin name followed by an underscore followed by the version. As an example, the filename and extension for the above plugin would be `myCustomPlugin_1.0.0.so`
* As documented in the previous step, plugins not utilizing versioning are simply named after the plugin name.