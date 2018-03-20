# How to Contribute to Kanali

## Getting Started
This project uses [dep](https://github.com/golang/dep) to manage dependencies. If not found locally, `make install` will bootstrap it locally.

Use the following commands to bootstrap a local development environment for Kanali. Note that Kanali is built using Go version `1.10`.

```sh
$ mkdir -p $GOPATH/src/github.com/northwesternmutual
$ cd $GOPATH/src/github.com/northwesternmutual
$ git clone https://github.com/northwesternmutual/kanali.git
$ cd kanali
$ make install
```

## Testing

This project contains a robust suite of unit and end-to-end (e2e) tests. Instructions detailing how to execute these suites are given below:

#### unit

```sh
$ make unit_test
```

#### e2e

This test suite requires a running Kubernetes cluster. Bootstrap a local one as follows:

```sh
$ minikube start --kubernetes-version v1.9.0 --feature-gates CustomResourceValidation=true
```
If you would like to run this test suite against local changes, it is necessary to create a container image with these changes. This can be done in the following manner. Note that many other aspects of this test suite are configurable via environment variables. Reference [this](./hack/e2e.sh) for a list of these variables.

```sh
$ eval $(minikube docker-env)
$ docker build -t kanali:local .
```

Execute the e2e test suite:

```sh
$ make e2e_test
```

## Project Structure

Below is an overview of this project's structure (not all files are listed):

```sh
github.com/northwesternmutual/kanali
  cmd/                - All binaries
    kanali/
      app/            - Code for binary
      main.go
  examples/           - Examples of all resources in the kanali.io API group
  hack/               - Miscellaneous project scripts
  docs/               - Hugo assets for documentation site
  helm/               - Helm chart for Kanali
  logo/               - Logo assets
  pkg/                - (See footnote 1)
  test/               - Integration and e2e tests
  Gopkg.toml          - Project dependencies
```

<sup>1</sup> `pkg` is a collection of utility packages used by the Kanali components without being specific to its internals. Utility packages are kept separate from the core codebase to keep it as small and concise as possible. If some utilities grow larger and their APIs stabilize, they may be moved to their own repository, to facilitate re-use by other projects.

## Imports Grouping
This projects adheres to the following pattern when grouping imports in Go files:
* imports from standard library
* imports from other projects
* imports from internal project

In addition, imports in each group must be sorted by length. For example:
```go
import (
  "context"
  "net/http"

  "go.uber.org/zap"
  opentracing "github.com/opentracing/opentracing-go"

  "github.com/northwesternmutual/kanali/pkg/tags"
  "github.com/northwesternmutual/kanali/pkg/log"
)
```

## Making a Change
Before making any significant changes, please [open an issue](https://github.com/northwesternmutual/kanali/issues). Discussing your proposed changes ahead of time will make the contribution process smooth for everyone.

Once we've discussed your changes and you've got your code ready, make sure that tests are passing (make test or make cover) and open your PR. Your pull request is most likely to be accepted if it:

* Includes tests for new functionality.
* Follows the guidelines in [Effective Go](https://golang.org/doc/effective_go.html) and the [Go team's common code review comments](https://github.com/golang/go/wiki/CodeReviewComments).
* Has a good commit message.

## License
By contributing your code, you agree to license your contribution under the terms of the [Apache License](https://github.com/jaegertracing/jaeger/blob/master/LICENSE).

If you are adding a new file it should have a header like below. The easiest way to add such header is to run `make fmt`.

```sh
// Copyright (c) 2018 Northwestern Mutual.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
```