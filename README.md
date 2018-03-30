<p align="center">
<img src="logo/logo_with_name.png" alt="Kanali" title="Kanali" width="50%" />
</p>

[![Travis](https://img.shields.io/travis/northwesternmutual/kanali/master.svg)](https://travis-ci.org/northwesternmutual/kanali)
[![Coveralls](https://img.shields.io/coveralls/northwesternmutual/kanali/master.svg)](https://coveralls.io/github/northwesternmutual/kanali)
[![OpenTracing Badge](https://img.shields.io/badge/OpenTracing-enabled-blue.svg)](http://opentracing.io)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/northwesternmutual/kanali)
[![Go Report Card](https://goreportcard.com/badge/github.com/northwesternmutual/kanali)](https://goreportcard.com/report/github.com/northwesternmutual/kanali)

Kanali is a lightweight, [Kubernetes](https://kubernetes.io/) native API management gateway that together with network policies provide a robust, open source solution to Kubernetes ingress, API management, and API security.

Notable features:

* **Kubernetes Native:** Kanali extends the Kubernetes API by using [Custom Resource Definitions](https://kubernetes.io/docs/concepts/api-extension/custom-resources/#customresourcedefinitions), allowing Kanali to be configured and used in the same way as native Kubernetes resources.
* **Performance Centric:** As a middleware component, Kanali is developed with performance as the highest priority!
* **Plugin Framework:** Need to perform complex transformations or integrations with a legacy system? Kanali provides a framework allowing developers to create, integrate, and version control custom plugins.
* **User-Defined Configurations:** Kanali gives you control over declaratively configuring how your proxy behaves. Need mutual TLS, dynamic service discovery, mock responses, etc.? Kanali makes it easy!
* **Robust API Management:** Fine grained API key authorization, JWT validation, quota policies, rate limiting, etc., these are some of the built in API management capabilities that Kanali provides.
* **Analytics & Monitoring:** Kanali integrates with [Prometheus](https://prometheus.io/) and [Grafana](https://grafana.com/) to provide a robust set of metrics and a customizable dashboard so that you can monitor the performance of your APIs in real time.
* **Production Ready:** [Northwestern Mutual](https://www.northwesternmutual.com/) uses Kanali in production at scale for all of its API management needs in their cloud native stack.
* **Cloud Native Deployment:** Kanali is deployed using native Kubernetes constructs which is assisted by an included [Helm](https://helm.sh/) chart.
* **Open Tracing Integration:** Kanali integrates with [OpenTracing](http://opentracing.io/), a vendor-neutral open standard for distributed tracing. [Jaeger](http://jaeger.readthedocs.io/en/latest/), a distributed tracing system, is supported out of the box to provide a visual representation for your traces.

## Getting Started

Find complete documentation at [kanali.io](https://kanali.io).

Try our [interactive tutorial](https://kanali.io/tutorial).

## Media

[KubeCon 2017](https://youtu.be/--LSmvCKVSs)

## Contributing
See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches and the contribution workflow.

## Support

Before [filing an issue](https://github.com/northwesternmutual/kanali/issues/new), please visit our [troubleshooting guide]().