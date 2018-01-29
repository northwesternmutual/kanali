# Plugin Guide

> a developer's guide to Kanali plugins

## Overview

Want to create you own plugin for *Kanali*? No problem! With *Kanali's* extensible plugin system, it's easy! Here are some steps to follow:

* [Step 1: Create](#step-1-create)
* [Step 2: Test](#step-2-test)
* [Step 3: Compile](#step-3-compile)
* [Step 4: Place](#step-4-place)
* [Step 5: Use](#step-5-use)
* [Step 6: Version](#step-6-optional-version)

## Step 1: Create

There are essentially only two main design requirements for a Kanali plugin:

1. The plugin must implement Kanali's [`Plugin`](https://github.com/northwesternmutual/kanali/blob/master/plugins/plugin.go) interface.
2. The name of the exported variable implementing the interface must be `Plugin`.

To make this simple, I've provided a [Yeoman](http://yeoman.io/) template that will dynamically scaffold everything out for you. You can find the template [here](https://github.com/northwesternmutual/kanali-plugin-template) or just follow the instructions below:

```sh
$ npm install -g yo
$ npm install -g generator-kanali-plugin
$ yo kanali-plugin
```

Every plugin has the ability to intercept a request at two points during the request lifecycle. It does this by invoking the two lifecycle methods that make up the [`Plugin`](https://github.com/northwesternmutual/kanali/blob/master/plugins/plugin.go) interface:

1. `OnRequest` is invoked *before* the request is proxied upstream.
2. `OnResponse` is invoked *after* the request has been returned from the upstream service but *before* the request is returned to the client.

If an error was encountered during the plugin logic that should result in the termination of the request, the method should return an [`error`](https://golang.org/pkg/errors/#pkg-examples). If no error was encountered and you would like the request's lifecycle to proceed as normal, return `nil`.

If you would like to relay a detailed error to a client which includes items like the status code, return [`*errors.Error`]() which is shown below. This utility type implements the `error` interface and allows for a verbose response to the client. If an `error` is returned that is not of this type, a default error will be relayed to the client.

```go
type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Details string `json:"details"`
}
```

A series of method parameters are provided for optional usage inside of your plugin logic. Here is a table detailing their purpose:

name        | type             | description
------------|------------------|----------------------|-----------|------------
`ctx`        | [`context.Context`](https://golang.org/pkg/context/)  | Parent [`opentracing.Span`](https://godoc.org/github.com/opentracing/opentracing-go#Span) is stored here.
`config`      | [`plugin.Config`]() | Arbitrary set of key value configuration items for use by a plugin.
`w`       | [`*httptest.ResponseRecorder`](https://golang.org/pkg/net/http/httptest/#ResponseRecorder) | Response details from upstream.
`r`        | [`*http.Request`](https://golang.org/pkg/net/http/#Request) | Original HTTP request.

## Step 2: Test

No code is complete without ample test coverage! If you are using the [template](https://github.com/northwesternmutual/kanali-plugin-template) to help bootstrap your plugin, your testing framework is already scaffolded for you. Simply run the following commands:

```sh
$ make test
$ make cover
```

## Step 3: Compile

Go plugins are not compiled into Kanali's binary but instead are loaded at runtime. Go expects are plugin to be an ELF shared object file. Here's an example of how to compile your plugin:

```sh
$ go build -buildmode=plugin -o myCustomPlugin.so myCustomPlugin.go
```

## Step 4: Place

During runtime, Kanali will look for compiled plugins in a location specified by the `plugins.location` cli flag with a default location of `/`.

In theory, plugins should be able to be dynamically loaded without the need to modify Kanali's container image. In the context of Kubernetes, the recommended approach would be to use [init containers](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/). The init container would move the compiled plugin into the location where the Kanali binary is configured to look for them.

**_However_**, there is a nasty bug that prevents this from being the case until at least Go v1.11. Evidence of this can be tracked in the following issues:
* [https://github.com/golang/go/issues/18827](https://github.com/golang/go/issues/18827)
* [https://github.com/golang/go/issues/20481](https://github.com/golang/go/issues/20481)

The root cause of the issue stems from the fact that unless a plugin is compiled using the exact same `vendor/` instance as Kanali, bad things will happen. For example, global variables will be duplicated, `init` functions will execute twice, etc.

The **_workaround_** , until this issue is fixed, entails building your custom plugin within the same `Dockerfile` as Kanali. The goal is to take the superset of dependencies that Kanali and every plugin requires and build them all using this aggregated dependency tree. This will only work if all intersecting dependencies utilize the exact same revision.

To simplify this as much as possible for the user, a `requirements.yaml` will be mounted into the container at build time that defines what plugins are to be included.

```yaml
plugins:
- name: apiKey
  source: github.com/northwesternmutual/kanali-plugin-apikey
  versions:
  - name: v2.0.0 # name of version to be used in proxy.Spec.Plugins[i].Version
    default: true # if true, this version will be used if proxy.Spec.Plugins[i].Version is not specified. Only one version can be specified as default.
    version: v2.0.0 # git version
  - name: abc123
    default: false
    revision: abc123 # git commit
  - name: master
    default: false
    branch: master # git branch
```
The resulting container image will have three different versions of the `apiKey` plugin ready for use.

## Step 5: Use

Adding a custom plugin to an `ApiProxy` is simple:

```yaml
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
  name: plugin-example
  namespace: default
spec:
  ...
  plugins:
  - name: apiKey
    config:
      key1: value1
      key2: value2
```

Each plugin allows a user to configure an arbitrary number of key value pairs to be passed to the plugin logic. Values are restricted to type `string`. An example use case might be a JWT plugin where the audience id associated with this `ApiProxy` is passed as a config value.

## Step 6 (optional): Version

Kanali gives you the option to version control your plugins. Below is an example of how to use a certain version of your plugin:

```yaml
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
  name: plugin-example
  namespace: default
spec:
  ...
  plugins:
  - name: apiKey
    version: v2.0.0
```

Here are some naming rules that must be followed to accomplish this (if you are following the workaround in step 4, the following steps are abstracted from the user):
* The file name for a compiled plugin representing a specific version must be the plugin name followed by an underscore followed by the version. As an example, the filename and extension for the above plugin would be `apiKey_v2.0.0.so`
* As documented in the previous step, plugins not utilizing versioning are simply named after the plugin name.
