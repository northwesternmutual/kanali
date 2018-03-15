<p align="center">
<img src="logo/logo_with_name.png" alt="Kanali" title="Kanali" width="50%" />
</p>

[![Travis](https://img.shields.io/travis/northwesternmutual/kanali/kanali.io%2Fv2alpha1.svg)](https://travis-ci.org/northwesternmutual/kanali)
[![Coveralls](https://img.shields.io/coveralls/northwesternmutual/kanali/kanali.io%2Fv2alpha1.svg)](https://coveralls.io/github/northwesternmutual/kanali)
[![OpenTracing Badge](https://img.shields.io/badge/OpenTracing-enabled-blue.svg)](http://opentracing.io)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/northwesternmutual/kanali)
[![Go Report Card](https://goreportcard.com/badge/github.com/northwesternmutual/kanali)](https://goreportcard.com/report/github.com/northwesternmutual/kanali)

Kanali is an efficient [Kubernetes](https://kubernetes.io/) ingress proxy with robust API management capabilities. Built using native Kubernetes constructs, Kanali gives you all the capabilities you need when exposing services in production without the need for multiple tools to accomplish them. Here are some notable features:

* **Kubernetes Native:** Kanali extends the Kubernetes API by using [Custom Resource Definitions](https://kubernetes.io/docs/concepts/api-extension/custom-resources/#customresourcedefinitions), allowing Kanali to be configured and used in the same way as native Kubernetes resources.
* **Performance Centric:** As a middleware component, Kanali is developed with performance as the highest priority! You could instantly improve your application's network performance by using Kanali.
* **Powerful, Decoupled Plugin Framework:** Need to perform complex transformations or integrations with a legacy system? Kanali provides a framework allowing developers to create, integrate, and version control custom plugins without every touching the Kanali codebase. Read more about plugins [here](https://github.com/northwesternmutual/kanali/blob/master/PLUGIN_GUIDE.md).
* **User-Defined Configurations:** Kanali gives you complete control over declaratively configuring how your proxy behaves. Need mutual TLS, dynamic service discovery, mock responses, etc.? No problem! Kanali makes it easy!
* **Robust API Management:** Fine grained API key authorization, quota policies, rate limiting, etc., these are some of the built in API management capabilities that Kanali provides. In addition, it follows native Kubernetes patterns for API key creation and binding making it easy and secure to control access to your proxy.
* **Analytics & Monitoring:** Kanali uses [Grafana](https://grafana.com/) and [Prometheus](https://prometheus.io/) to provide a customizable and visually appealing experience so that you can get real time alerting and visualization around Kanali's metrics. Find out more [here](#analytics-and-monitoring)!
* **Production Ready:** [Northwestern Mutual](https://www.northwesternmutual.com/) uses Kanali in Production to proxy, manage, and secure all Kubernetes hosted services.
* **Easy Installation:** Kanali does not rely on an external database, infrastructure agents or workers, dedicated servers, etc. Kanali is deployed in the same manner as any other service in Kubernetes. Find installation instructions [here](#installation)
* **Open Tracing Integration:** Kanali integrates with [Open Tracing](http://opentracing.io/), hosted by the [Cloud Native Foundation](https://www.cncf.io/), which provides consistent, expressive, vendor-neutral APIs allowing you to trace the entire lifecycle of a request. [Jaeger](http://jaeger.readthedocs.io/en/latest/), a distributed tracing system open sourced by Uber Technologies and recently accepted into the Cloud Native Foundation, is supported out of the box to provide a visual representation for your traces.

## Getting Started

Find complete documentation at [kanali.io](https://kanali.io).

Try our [interactive tutorial](https://kanali.io/tutorial).

## Contributing
See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches and the contribution workflow.

## Support

Before filing an issue, make sure you visit our [troubleshooting guide]().