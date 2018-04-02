+++
description = "Configuration flags"
title = "Flags"
date = "2017-04-10T16:43:08+01:00"
draft = false
weight = 100
bref="Configuration flags"
toc = true
+++

### Usage

Each of the following flags can be set to a custom value. Values can also be set via environment variables and configuration files. The name of the corresponding environment variable is documented below. The order of precedence is as follows (each item has precedence over the item below it):

<span class="label focus">Flag</span>
<br />
<span class="label success">Environment Variable</span>
<br />
<span class="label warning">Config file</span>

To illustrate an example, let's assume we have a `config.toml` with the following value:
```
[kubernetes]
kubeconfig = "/usr/frank/.kube/config"
```

If the `KANALI_KUBERNETES_KUBECONFIG` environment variable is set, it will tag precedence over the value in the `config.toml` file.

If Kanali is started with the `--kubernetes.kubeconfig` flag, that will take the ultimate precedence over the previous two options.

### Configuration Flags

#### `--kubernetes.kubeconfig`
**Description:** Location of Kubernetes config file. This flag is most likely only applicable if running Kanali outside of a Kubernetes cluster.

**Type:** *string*

**Environment Variable:** `KANALI_KUBERNETES_KUBECONFIG`

<br />

#### `--plugins.location`
**Description:** Location of directory containing compiled shared object plugin files.

**Type:** *string*

**Default:** `/`

**Environment Variable:** `KANALI_PLUGINS_LOCATION`

<br />

#### `--plugins.apiKey.decryption_key_file`
**Description:** Location of PEM encoded RSA private key.

**Type:** *string*

**Environment Variable:** `KANALI_PLUGINS_APIKEY_DECRYPTION_KEY_FILE`

<br />

#### `--plugins.apiKey.header_key`
**Description:** Name of HTTP header containing API key value.

**Type:** *string*

**Default:** `apikey`

**Environment Variable:** `KANALI_PLUGINS_APIKEY_HEADER_KEY`

<br />

#### `--process.log_level`
**Description:** Kanali logging level

**Type:** *string*

**Default:** `info`

**Environment Variable:** `KANALI_PROCESS_LOG_LEVEL`

<br />

#### `--profiling.insecure_port`
**Description:** TCP port value that the profiling server will bind to for an HTTP server.

**Type:** *number*

**Default:** `9090`

**Environment Variable:** `KANALI_PROFILING_INSECURE_PORT`

<br />

#### `--profiling.insecure_bind_address`
**Description:** Bind address that the profiling server will bind to for an HTTP server.

**Type:** *string*

**Default:** `0.0.0.0`

**Environment Variable:** `KANALI_PROFILING_INSECURE_BIND_ADDRESS`

<br />

#### `--prometheus.secure_port`
**Description:** TCP port value that the profiling server will bind to for an HTTPS server.

**Type:** *number*

**Default:** `0`

**Environment Variable:** `KANALI_PROMETHEUS_SECURE_PORT`

<br />

#### `--prometheus.insecure_port`
**Description:** TCP port value that the profiling server will bind to for an HTTP server.

**Type:** *number*

**Default:** `9000`

**Environment Variable:** `KANALI_PROMETHEUS_INSECURE_PORT`

<br />

#### `--prometheus.insecure_bind_address`
**Description:** Bind address that the profiling server will bind to for an HTTP server.

**Type:** *string*

**Default:** `0.0.0.0`

**Environment Variable:** `KANALI_PROMETHEUS_INSECURE_BIND_ADDRESS`

<br />

#### `--prometheus.secure_bind_address`
**Description:** Bind address that the profiling server will bind to for an HTTPS server.

**Type:** *string*

**Default:** `0.0.0.0`

**Environment Variable:** `KANALI_PROMETHEUS_SECURE_BIND_ADDRESS`

<br />

#### `--proxy.enable_cluster_ip`
**Description:** Whether or not to use the Kubernetes `clusterIP` instead of dns for Kubernetes upstream services.

**Type:** *boolean*

**Default:** `true`

**Environment Variable:** `KANALI_PROXY_ENABLE_CLUSTER_IP`

<br />

#### `--proxy.header_mask_Value`
**Description:** When logging and tracing a request's headers, this flag specifies the value to use when masking certain header's values.

**Type:** *string*

**Default:** `omitted`

**Environment Variable:** `KANALI_PROXY_HEADER_MASK_VALUE`

<br />

#### `--proxy.enable_mock_responses`
**Description:** Whether or not to enable mock responses.

**Type:** *boolean*

**Default:** `true`

**Environment Variable:** `KANALI_PROXY_ENABLE_MOCK_RESPONSES`

<br />

#### `--proxy.upstream_timeout`
**Description:** Upstream timeout duration.

**Type:** *duration*

**Default:** `0h0m10s`

**Environment Variable:** `KANALI_PROXY_UPSTREAM_TIMEOUT`

<br />

#### `--proxy.mask_header_keys`
**Description:** When logging and tracing a request's headers, this flag specifies the headers (comma separated) whose values will be masked.

**Type:** *string*

**Environment Variable:** `KANALI_PROXY_MASK_HEADER_KEYS`

<br />

#### `--proxy.tls_common_name_validation`
**Description:** This flag specifies whether server name validation will be performed during an upstream TLS handshake.

**Type:** *boolean*

**Default:** `true`

**Environment Variable:** `KANALI_PROXY_TLS_COMMON_NAME_VALIDATION`

<br />

#### `--proxy.default_header_values`
**Description:** Default values for HTTP headers when using dynamic service discovery. Flag value takes the form `foo=bar,car=baz`.

**Type:** *string*

**Environment Variable:** `KANALI_PROXY_DEFAULT_HEADER_VALUES`

<br />

#### `--server.secure_port`
**Description:** TCP port value that the profiling server will bind to for an HTTPS server.

**Type:** *number*

**Default:** `0`

**Environment Variable:** `KANALI_SERVER_SECURE_PORT`

<br />

#### `--server.insecure_port`
**Description:** TCP port value that the profiling server will bind to for an HTTP server.

**Type:** *number*

**Default:** `8080`

**Environment Variable:** `KANALI_SERVER_INSECURE_PORT`

<br />

#### `--server.insecure_bind_address`
**Description:** Bind address that the profiling server will bind to for an HTTP server.

**Type:** *string*

**Default:** `0.0.0.0`

**Environment Variable:** `KANALI_SERVER_INSECURE_BIND_ADDRESS`

<br />

#### `--server.secure_bind_address`
**Description:** Bind address that the profiling server will bind to for an HTTPS server.

**Type:** *string*

**Default:** `0.0.0.0`

**Environment Variable:** `KANALI_SERVER_SECURE_BIND_ADDRESS`

<br />

#### `--server.tls.cert_file`
**Description:** Location of the PEM encoded TLS certificate for the Kanali gateway server.

**Type:** *string*

**Environment Variable:** `KANALI_SERVER_TLS_CERT_FILE`

<br />

#### `--server.tls.key_file`
**Description:** Location of the PEM encoded TLS private key for the Kanali gateway server.

**Type:** *string*

**Environment Variable:** `KANALI_SERVER_TLS_KEY_FILE`

<br />

#### `--server.tls.ca_file`
**Description:** Location of the PEM encoded TLS certificate authority for the Kanali gateway server. This is specified when validating the identity of the client.

**Type:** *string*

**Environment Variable:** `KANALI_SERVER_TLS_CA_FILE`

<br />

#### `--tracing.config_file`
**Description:** Location of the Jaeger config file in `yaml` format. Configuration details can be found [here](https://godoc.org/github.com/uber/jaeger-client-go/config#Configuration).

**Type:** *string*

**Environment Variable:** `KANALI_TRACING_CONFIG_FILE`

<br />
