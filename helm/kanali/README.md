# Kanali

> A Kubernetes Native Ingress Proxy and API Management Solution.

## Installing the Chart

To install the chart with the release name `my-release`:

```sh
$ git clone https://github.com/northwesternmutual/kanali.git
$ kubectl apply -f kanali/hack/helm-rbac.yaml
$ helm init --service-account tiller
$ helm install --name my-release ./kanali/helm/kanali
```

## Uninstalling the Chart

To uninstall/delete the my-release deployment:

```sh
$ helm del --purge my-release
```

## Configuration

| Parameter                                 | Description                                          | Default                   |
| ----------------------------------------- | ---------------------------------------------------- | ------------------------- |
| `kanali.namespace`                        | Namespace for Kanali resources                       | default                   |
| `kanali.logLevel`                         | Log level                                            | info                      |
| `kanali.gateway.securePort`               | HTTPS server port                                    | 8443                      |
| `kanali.gateway.secureBindAddress`        | HTTP server bind address                             | 0.0.0.0                   |
| `kanali.gateway.insecurePort`             | HTTP server port                                     | 0 (disabled)              |
| `kanali.gateway.insecureBindAddress`      | HTTP server bind address                             | 0.0.0.0                   |
| `kanali.gateway.rsa.secretName`           | Secret containing private key for API key decryption | kanali-rsa                |
| `kanali.gateway.tls.secretName`           | Secret containing tls assets for HTTPS server        | kanali-pki                |
| `kanali.gateway.tls.verifyClient`         | Should mutual TLS be used for HTTPS server           | false                     |
| `kanali.gateway.image.pullPolicy`         | Container pull policy                                | IfNotPresent              |
| `kanali.gateway.image.image`              | Container image                                      | northwesternmutual/kanali |
| `kanali.gateway.image.tag`                | Container image tag                                  | local                     |
| `kanali.gateway.scale.minReplicas`        | Minimum number of Kanali replicas                    | 1                         |
| `kanali.gateway.scale.maxReplicas`        | Maximum number of Kanali replicas                    | 5                         |
| `kanali.gateway.scale.targetCPU`          | Target CPU percentage before a scaling event         | 500                       |
| `kanali.gateway.plugins.apikey.headerKey` | HTTP header key for apikey value                     | apikey                    |
| `kanali.jaeger.config`                    | Jaeger YAML config                                   | *see values.yaml*         |
| `kanali.prometheus.securePort`            | HTTPS server port                                    | 0 (disabled)              |
| `kanali.prometheus.secureBindAddress`     | HTTP server bind address                             | 0.0.0.0                   |
| `kanali.prometheus.insecurePort`          | HTTP server port                                     | 9000                      |
| `kanali.prometheus.insecureBindAddress`   | HTTP server bind address                             | 0.0.0.0                   |
| `kanali.profiler.insecurePort`            | HTTP server port                                     | 9090                      |
| `kanali.profiler.insecureBindAddress`     | HTTP server bind address                             | 0.0.0.0                   |