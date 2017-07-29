# ApiProxy

| Field | Required | Description |
| ----- | -------- | ----------- |
| apiVersion<br />*string*   | `true`       |   APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values.   |
| kind<br />*string*   | `true`      |    Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase.         |
| metadata<br />*[ObjectMeta](https://kubernetes.io/docs/api-reference/v1.6/#objectmeta-v1-meta)*  | `true`    |     Standard object's metadata.        |
| spec<br />*[ApiProxySpec](#apiproxyspec)*   | `true`     |      Defines an ApiProxy   |

# ApiProxySpec

Field | Required | Description |
| ----- | -------- | ----------- |
| path<br />*string*   | `true`       |   Declares what incoming request to be correlated to this proxy (must be unique although subsets are allowed). Must start with a `/`.   |
| target<br />*string*   | `false`      |    Declares the first beginning subset of the upstream path. The complement of the the incoming path and the proxy path will be concatenated onto the end of the target path. Must start with a `/`.         |
| mock<br />[*Mock*](#mock)   | `false`      |    if mock if defined and *Kanali* is started with the `--mock-enabled` flag, the mock responses will be used instead of proxying to the actual backend service.         |
| hosts<br />*[Host](#host) array*  | `false`    |     Specifies what destination host(s) to match against when using SNI.        |
| service<br />[*Service*](#service)   | `true`     |     Specifies how to discover a Kubernetes service. *NOTE:* to comply with Kubernetes conventions, the namespace of the service will match the namespace of the ApiProxy        |
| plugins<br />*[Plugin](#plugin) array*   | `false`      |    Specifies what plugins, if any, to use throughout the request's lifecycle. All plugins have the opportunity to intercept a request both before and after the proxy pass.         |
| ssl<br />[*SSL*](#ssl)   | `false`       |      Specifies the details of the TLS connection to configure for the upstream request. *NOTE:* this SSL object is overridden if SNI is used. If a host is specified and SNI is not used, this SSL object takes precedence for that specific upstream.       |

# Mock

| Field | Required | Description |
| ----- | -------- | ----------- |
| configMapName<br />*string*  | `true` | Name of the Kubernetes ConfigMap that configures the mock responses. Must live in the same namespace as this Proxy  |

# Host

| Field | Required | Description |
| ----- | -------- | ----------- |
| name<br />*string*   | `true`       |   Name of the destination host to use for SNI.   |
| ssl<br />[*SSL*](#ssl)   | `true`       |      Specifies the details of the TLS connection to configure for this host.    |

# Service

| Field | Required | Description |
| ----- | -------- | ----------- |
| name<br />*string*  | If undefined, *labels* must be defined. | Name of the Kubernetes service to discover. *NOTE:* to comply with Kubernetes conventions, the namespace of the service will match the namespace of the ApiProxy.   |
| port<br />*int*   | `true`       |   The http port to use.   |
| labels<br />*[Label](#label) array*   | If undefined, *name* must be defined.       |   List of labels to use to discover Kubernetes services. If multiple found, fist found will be used.   |

# Label

| Field | Required | Description |
| ----- | -------- | ----------- |
| name<br />*string*  | `true` | Name of the label on the Kubernetes service.  |
| value<br />*string*  | If undefined, *header* must be specified. | Value of the Kubernetes service label corresponding to the label with the above name.  |
| header<br />*string*  | If undefined, *header* must be specified. | Name of the http header whose value will be matched against the value of the label corresponding to the label with the above name.  |

# Plugin

| Field | Required | Description |
| ----- | -------- | ----------- |
| name<br />*string*   | `true`       |   Name of the plugin   |
| version<br />*string*    | `false`       |      Specifies the details of the TLS connection to configure for this host.    |

# SSL

| Field | Required | Description |
| ----- | -------- | ----------- |
| secretName<br />*string*  | `true` | Name of the Kubernetes secret to use. *NOTE:* Secret type **must be** *kubernetes.io/tls*  |

