# Kanali

[![Travis](https://img.shields.io/travis/northwesternmutual/kanali/master.svg)](https://travis-ci.org/northwesternmutual/kanali) [![Coveralls](https://img.shields.io/coveralls/northwesternmutual/kanali/master.svg)](https://coveralls.io/github/northwesternmutual/kanali) [![Documentation](https://img.shields.io/badge/docs-latest-brightgreen.svg)](https://github.com/northwesternmutual/kanali/tree/master/docs/docs.md) [![OpenTracing Badge](https://img.shields.io/badge/OpenTracing-enabled-blue.svg)](http://opentracing.io) [![Tutorial](https://img.shields.io/badge/tutorial-postman-orange.svg)](http://tutorial.kanali.io) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/northwesternmutual/kanali)
[![Go Report Card](https://goreportcard.com/badge/github.com/northwesternmutual/kanali)](https://goreportcard.com/report/github.com/northwesternmutual/kanali)

Kanali is an extremely efficient [Kubernetes](https://kubernetes.io/) ingress controller with robust API management capabilities. Built using native Kubernetes constructs, Kanali gives you all the capabilities you need when exposing services in production without the need for multiple tools to accomplish them. Here are some notable features:

* **Kubernetes Native:** Kanali extends the Kubernetes API by using [Third Party Resources](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-third-party-resource/), allowing Kanali to be configured in the same way as native Kubernetes resources.
* **Performance Centric:** As a middleware component, Kanali is developed with performance as the highest priority! You could instantly improve your application's network performance by using Kanali.
* **Powerful, Decoupled Plugin Framework:** Need to perform complex transformations or integrations with a legacy system? Kanali provides a framework allowing developers to create, integrate, and version control custom plugins without every touching the Kanali codebase. Read more about plugins [here](https://github.com/northwesternmutual/kanali/blob/master/PLUGIN_GUIDE.md).
* **User-Defined Configurations:** Kanali gives you complete control over declaratively configuring how your proxy behaves. Need mutual TLS, dynamic service discovery, mock responses, etc.? No problem! Kanali makes it easy!
* **Robust API Management:** Fine grained API key authorization, quota policies, rate limiting, etc., these are some of the built in API management capabilities that Kanali provides. In addition, it follows native Kubernetes patterns for API key creation and binding making it easy and secure to control access to your proxy.
* **Analytics & Monitoring:** Kanali uses [Grafana](https://grafana.com/) and [InfluxDB](https://www.influxdata.com/) to provide a customizable and visually appealing experience so that you can get real time alerting and visualization around Kanali's metrics. Find out more [here](#analytics-and-monitoring)!
* **Production Ready:** [Northwestern Mutual](https://www.northwesternmutual.com/) uses Kanali in Production to proxy, manage, and secure all Kubernetes hosted services.
* **Easy Installation:** Kanali does not rely on an external database, infrastructure agents or workers, dedicated servers, etc. Kanali is deployed in the same manner as any other service in Kubernetes. Find installation instructions [here](#installation)
* **Open Tracing Integration:** Kanali integrates with [Open Tracing](http://opentracing.io/), endorsed by the [Cloud Native Foundation](https://www.cncf.io/), which provides consistent, expressive, vendor-neutral APIs allowing you to trace the entire lifecycle of a request. [Jaeger](http://jaeger.readthedocs.io/en/latest/), a distributed tracing system open sourced by Uber Technologies, is supported out of the box providing a visual representation for your traces.

# Table of Contents

* [Quick Start](#quick-start)
* [Tutorial](#tutorial)
* [Documentation](#documentation)
* [Plugins](#plugins)
* [Analytics, Monitoring, and Tracing](#analytics-monitoring-and-tracing)
* [Installation](#installation)
  * [Helm](#helm)
  * [Manual](#manual)
* [Local Development](#local-development)
* [Usage and Configuration](#usage-and-configuration)
  * [CLI Flags](#cli-flags)
  * [Environment Variables](#environment-variables)
  * [Configuration Files](#configuration-files)

# Quick start

```sh
$ git clone git@github.com:northwesternmutual/kanali.git && cd kanali
$ minikube start
$ ./scripts/install.sh # wait until all pods are in running state
$ kubectl apply -f ./examples/exampleOne.yaml
$ curl $(minikube service kanali --url --format="https://{{.IP}}:{{.Port}}")/api/v1/example-one
$ open $(minikube service kanali-grafana --url)/dashboard/file/kanali.json
$ open $(minikube service jaeger-all-in-one --url)
```

# Tutorial

A complete guide showcasing all of Kanali's features is provided to ease the learning curve! The guided tutorial can be found [here](http://tutorial.kanali.io).

# Documentation

Looking for documentation for the custom Kubernetes resources that Kanali creates and uses? Find it [here](./docs/docs.md).

# Plugins

While Kanali has its own built in plugins, it boasts a decoupled plugin framework that makes it easy for anyone to write and integrate their own custom and version controlled plugin! The guided tutorial can be found [here](./PLUGIN_GUIDE.md).

# Analytics, Monitoring, and Tracing

Kanali leverages [Grafana](https://grafana.com/) and [InfluxDB](https://www.influxdata.com/) for analytics and monitoring. It also uses [Jaeger](http://jaeger.readthedocs.io/en/latest/) for tracing. If you are using Helm to deploy Kanali, these tools are deployed and configured for you.

Jaeger                                           | Grafana             
:-----------------------------------------------:|:-------------------------:
<img src="./assets/jaeger1.png" width="600">         | <img src="./assets/grafana.png" width="600">
<img src="./assets/jaeger2.png" width="600">  |  

# Installation

There are multiple ways to deploy Kanali. For each, a Kubernetes cluster is required. For local testing and development, use [Minikube](https://github.com/kubernetes/minikube) to bootstrap a cluster locally.

Each option below requires cloning this project locally. You can do so with the following command:

```sh
$ git clone https://github.com/northwesternmutual/kanali.git
# $ git clone git@github.com:northwesternmutual/kanali.git
$ cd kanali
```

## Helm

[Helm](https://github.com/kubernetes/helm) is a tool for managing Kubernetes charts. Charts are packages of pre-configured Kubernetes resources. Install Kanali along with Grafana, Influxdb, and Jaeger by using the following commands:

```sh
$ which helm || curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh | bash
$ helm init
$ ./scripts/helm-permissions.sh
$ kubectl get pods # wait until all pods are up and running
$ helm install ./helm --name kanali
```

## Manual

This installation process only installs Kanali. The deployment of Grafana, Influxdb, and Jaeger as left to the user.

```sh
$ which kubectl || curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/darwin/amd64/kubectl && chmod +x kubectl && mv kubectl /usr/bin/kubectl
$ kubectl apply -f kubernetes.yaml
```

# Local development

Below are the steps to follow if you want to build/run/test locally. [Glide](https://glide.sh/) is a dependency.

```sh
$ mkdir -p $GOPATH/src/github.com/northwesternmutual
$ cd $GOPATH/src/github.com/northwesternmutual
$ git clone https://github.com/northwesternmutual/kanali
$ cd kanali
$ make install_ci
$ make kanali
$ ./kanali --help
```

# Usage and Configuration

```sh
$ kanali [command] [flags]
```

```sh
start
    -y, --apikey-header-key string           Name of the HTTP header holding the apikey. (default "apikey")
    -b, --bind-address string                Network address that Kanali will listen on for incoming requests. (default "0.0.0.0")
    -d, --decryption-key-file string         Path to valid PEM-encoded private key that matches the public key used to encrypt API keys.
    -v, --disable-tls-cn-validation          Disable common name validate as part of an SSL handshake.
    -u, --enable-cluster-ip                  Enables to use of cluster ip as opposed to Kubernetes DNS for upstream routing.
    -m, --enable-mock                        Enables Kanali's mock responses feature. Read the documentation for more information.
    -t, --enable-proxy-protocol              Maintain the integrity of the remote client IP address when incoming traffic to Kanali includes the Proxy Protocol header.
    -f, --header-mask-value string           Sets the value to be used when omitting header values. (default "ommitted")
    -h, --help                               help for start
    -i, --influxdb-addr string               Influxdb address. Addr should be of the form 'http://host:port' or 'http://[ipv6-host%zone]:port' (default "monitoring-influxdb.kube-system.svc.cluster.local")
    -s, --influxdb-database string           Influxdb database (default "k8s")
    -r, --influxdb-password string           Influxdb password
    -q, --influxdb-username string           Influxdb username
    -n, --jaeger-agent-url string            Endpoint to the Jaeger agent (default "jaeger-all-in-one-agent.default.svc.cluster.local")
    -j, --jaeger-sampler-server-url string   Endpoint to the Jaeger sampler server (default "jaeger-all-in-one-agent.default.svc.cluster.local")
    -p, --kanali-port int                    Sets the port that Kanali will listen on for incoming requests.
    -l, --log-level string                   Sets the logging level. Choose between 'debug', 'info', 'warn', 'error', 'fatal'. (default "info")
    -o, --peer-udp-port int                  Sets the port that all Kanali instances will communicate to each other over. (default 10001)
    -g, --plugins-location string            Location of custom plugins shared object (.so) files. (default "/")
    -a, --tls-ca-file string                 Path to x509 certificate authority bundle for mutual TLS.
    -c, --tls-cert-file string               Path to x509 certificate for HTTPS servers.
    -k, --tls-private-key-file string        Path to x509 private key matching --tls-cert-file.
    -w, --upstream-timeout string            Set length of upstream timeout. Defaults to none (default "0h0m0s")

version
    no flags

help
    no flags
```

## CLI Flags

See above for all the cli flags available.

## Environment Variables

For each flag above, there is a corresponding environment variables. This variable name mirrors the cli flag name with the following modifications:
* Every character is upper case,
* The `-` character is replaced by the `_` character
* Everything is prefixed with `KANALI_`

As an example, if I wanted to overwrite the `--bind-address` flag, I would set the `KANALI_BIND_ADDRESS` environment variable.

## Configuration Files

A flag can also be overwritten via configuration files. Kanali accepts configuration files in `JSON`, `YAML`, `TOML`, or `HCL` formats.

Kanali will look for files in the following locations (in order of precedence):
* `/etc/kanali/config.ext`
* `$HOME/conifig.ext`
* `./conifig.ext`

*NOTE*: For nested objects in configuration files, the corresponding cli flag or environment variable includes a `.` in between every level. For example, a cli flag of `--foo.bar=car` would correspond to the following JSON:

```json
{
  "foo": {
    "bar": "car"
  }
}
```
