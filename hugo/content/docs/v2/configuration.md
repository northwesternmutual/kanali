+++
date = "2017-04-10T16:41:54+01:00"
weight = 20
description = "Learn how to declaratively configure Kanali"
title = "Configuration"
draft = false
bref= "Learn how to declaratively configure your API"
toc = true
+++

### Introduction
The goal of this section is to provide an introduction to each of Kanali's configurable resources. More in depth documentation for each resource can be found [here](https://godoc.org/github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2).

<br/>

### The `ApiProxy` Resource

This resource declaratively defines how your upstream services are reached. There are 3 main sections that can be configured for this resource.

<span class="label focus">source</span> how your proxy will be accessible

<span class="label focus">target</span> how your upstream service will be reached

<span class="label focus">plugins</span> how your proxy will be secured

Let's explore each in more detail.

#### Source

There are just 2 ways to configure a source, with or without a virtual host. If <code>virtualHost</code> is specified than, that proxy will only be discoverable if that host is being used. If <code>virtualHost</code> is not specified, than that proxy will be discoverable for all hosts. Here are examples of each.

<div class="example">
  <nav id="livetabs" data-component="tabs" data-live=".tab-live1"></nav>
  <div id="tab-basic" data-title="Basic" class="tab-live1">
  <pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
  name: example
spec:
  source:
    path: /foo
  target:
    path: /
    backend:
      endpoint: https://foo.bar.com:8443</pre>
</div>

<div id="tab-vhost" data-title="Virtual Host" class="tab-live1">
<pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
  name: example
spec:
  source:
    path: /bar
    virtualHost: foo.bar.com
  target:
    path: /
    backend:
      endpoint: https://foo.bar.com:8443</pre>
</div>
</div>

#### Target

There are 3 main sections that can be configured for this resource.

<span class="label focus">path</span> upstream path

<span class="label focus">backend</span> backend type

<span class="label focus">ssl</span> tls configuration

##### Path

This fields specifies the upstream path that an upstream request will be proxied to. In the below example, if a request is made to <code>/foo/baz</code>, the upstream service will see <code>/bar/baz</code>.

<div class="example">
<pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
  name: example
spec:
  source:
    path: /foo
  target:
    path: /bar
    backend:
      endpoint: https://foo.bar.com:8443</pre>
</div>

##### Backend

There are 4 different types of backends that you can configure. Toggle through each type below to learn more.

<div class="example">
  <nav id="livetabs" data-component="tabs" data-live=".tab-live"></nav>

  <div id="tab-service-static" data-title="Service (static)" class="tab-live">
    <pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
  name: example
spec:
  source:
    path: /example
  target:
    path: /
    backend:
      service:
        name: serviceName
        port: 8080</pre>

Statically defining an upstream Kubernetes service is the easiest way define an <code>ApiProxy</code>. Simply specify the name of your upstream Kubernetes service and the port your service is listening on.

<br />
<br />

The namespace of this <code>ApiProxy</code> must match the namespace of the upstream service.

  </div>
  <div id="tab-service-dynamic" data-title="Service (dynamic)" class="tab-live">
    <pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
 name: example
spec:
 source:
   path: /example
 target:
   path: /
   backend:
     service:
       port: 8080
       labels:
       - name: key
         value: value
       - name: deploy
         header: x-foo-deployment</pre>

Dynamically defining an upstream Kubernetes service provides an easy way to dynamically route traffic. To configure dynamic service discovery, simply specify the port your service is listening on and a set of labels.

<br />
<br />

Labels work in a similar fashion to that of Kubernetes match labels. The <code>name</code> field of each label will be matched against a metadata label name on Kubernetes services. The second field of each label can either be a <code>value</code> or <code>header</code>. If <code>value</code> is specified, it corresponds directly to the Kubernetes service metadata label value. If <code>header</code> is specified, than the value of the Kubernetes service label will be matched against the value of the HTTP header specified by this label. Let's look at a quick example.

<br />
<br />

Using the above <code>ApiProxy</code> as an example, If I make an request to <code>/example</code> and include the header <code>x-foo-deployment: bar</code>, Kanali will look for services in the <code>default</code> namespace that have at least the two following metadata labels.

<br />
<br />

<pre>
metadata:
 labels:
   key: value
   deploy: bar
</pre>

If multiple upstream services are found, Kanali will use the first one and add the response header <code>x-kanali-service-cardinality</code> whose value matches the cardinality of the discovered upstream services to aid in troubleshooting.

<br />
<br />

The namespace of this <code>ApiProxy</code> must match the namespace of the upstream service.

  </div>
  <div id="tab-endpoint" data-title="Endpoint" class="tab-live">
  <pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
 name: example
spec:
 source:
   path: /example
 target:
   path: /
   backend:
     endpoint: https://foo.bar.com:8443</pre>

There may be times when you want to add API management to an upstream service that is not deployed to Kubernetes. To prevent you from deploying a different API management gateway to solve this specific use case, Kanali lets you proxy to arbitrary endpoints. To configure an arbitrary endpoint, just specify it as shown above. The value must have the following structure.

<br />
<br />

<code>&lt;scheme&gt;://&lt;host&gt;</code>

<br />
<br />

The scheme must either be <code>http</code> or <code>https</code> and if the host does not contain a port, <code>80</code> will be used as the default.

</div>
  <div id="tab-mock" data-title="Mock" class="tab-live">
  <pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
 name: example
spec:
 source:
   path: /example
 target:
   path: /
   backend:
     mock:
       mockTargetName: mockTargetName</pre>

If this upstream type is used, Kanali will not make any request to any upstream service. Instead, it will return a preconfigured response. This response is defined in the <code>MockTarget</code> resource. This resource is explained in detail below.

<br />
<br />

The namespace of this <code>ApiProxy</code> must match the namespace of the <code>MockTarget</code> resource.

</div>
</div>

##### SSL

The presence of the <code>ssl</code> field specifies that tls will be used to secure the connection between Kanali and an upstream service. To configure this option, just specify the secret name containing the tls assets. An example is demonstrated below.

<div class="example">
<pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
  name: example
spec:
  source:
    path: /foo
  target:
    path: /bar
    backend:
      endpoint: https://foo.bar.com:8443
    ssl:
      secretName: my-secret</pre>
</div>

Let's assume that the specified secret above is structured like the example below. Note the presence of the <code>kanali.io/enabled</code> annotation. This annotation declares that Kanali is <i>allowed</i> to use this secret (this is due to Kubernetes RBAC limitations).

Note the data fields present in this secret. If your upstream service wants to perform client side validation, the tls certificate/key pair as specified in the <code>tls.crt</code> and <code>tls.key</code> fields will be send to the server.

<div class="example">
<pre>
---
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  annotations:
    kanali.io/enabled: 'true'
type: Opaque
data:
  tls.crt: <tls crt data>
  tls.key: <tls key data>
  tls.ca: <tls ca data></pre>
</div>

If you want to customize the name of the data keys, you can specify your custom key via an annotation. For example, if you want to use the data key <code>crt.pem</code> instead of <code>tls.crt</code>, you would need to include the annotation <code>kanali.io/cert: 'crt.pem'</code>. A complete list of override annotations for the data fields are listed below.

<table>
  <tr><td>Data field</td><td>Annotation</td></tr>
  <tr><td><code>tls.ca</code></td><td><code>kanali.io/ca: 'custom.ca.value'</code></td></tr>
  <tr><td><code>tls.crt</code></td><td><code>kanali.io/cert: 'custom.cert.value'</code></td></tr>
  <tr><td><code>tls.key</code></td><td><code>kanali.io/key: 'custom.key.value'</code></td></tr>
<table>

#### Plugins

Plugins enable the execution of encapsulated logic on a per proxy basis. Plugins are configured as a list of different plugins that you want executed for a specific <code>ApiProxy</code>. Each plugin in the list requires the name of the plugin and an optional config field specifying proxy level configuration items that will be passed to the plugin upon execution.

For a complete list of available plugins and their corresponding documentation, visit the [documentation for plugins](/docs/v2/plugins).

<div class="example">
<pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
  name: example
spec:
  source:
    path: /foo
  target:
    path: /bar
    backend:
      endpoint: https://foo.bar.com:8443
  plugins:
  - name: apikey
    config:
      bindingName: my-binding
  - name: jwt
    config:
      audienceID: abc123</pre>

</div>

##### SSL

The presence of the <code>ssl</code> field specifies that tls will be used to secure the connection between Kanali and an upstream service. To configure this option, just specify the secret name containing the tls assets. An example is demonstrated below.

<div class="example">
<pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
  name: example
spec:
  source:
    path: /foo
  target:
    path: /bar
    backend:
      endpoint: https://foo.bar.com:8443
    ssl:
      secretName: my-secret</pre>
</div>

Let's assume that the specified secret above is structured like the example below. Note the presence of the <code>kanali.io/enabled</code> annotation. This annotation declares that Kanali is <i>allowed</i> to use this secret (this is due to Kubernetes RBAC limitations).

Note the data fields present in this secret. If your upstream service wants to perform client side validation, the tls certificate/key pair as specified in the <code>tls.crt</code> and <code>tls.key</code> fields will be send to the server.

<div class="example">
<pre>
---
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  annotations:
    kanali.io/enabled: 'true'
type: Opaque
data:
  tls.crt: <tls crt data>
  tls.key: <tls key data>
  tls.ca: <tls ca data></pre>
</div>

If you want to customize the name of the data keys, you can specify your custom key via an annotation. For example, if you want to use the data key <code>crt.pem</code> instead of <code>tls.crt</code>, you would need to include the annotation <code>kanali.io/cert: 'crt.pem'</code>. A complete list of override annotations for the data fields are listed below.

<table>
  <tr><td>Data field</td><td>Annotation</td></tr>
  <tr><td><code>tls.ca</code></td><td><code>kanali.io/ca: 'custom.ca.value'</code></td></tr>
  <tr><td><code>tls.crt</code></td><td><code>kanali.io/cert: 'custom.cert.value'</code></td></tr>
  <tr><td><code>tls.key</code></td><td><code>kanali.io/key: 'custom.key.value'</code></td></tr>
<table>

#### Plugins

Plugins enable the execution of encapsulated logic on a per proxy basis. Plugins are configured as a list of different plugins that you want executed for a specific <code>ApiProxy</code>. Each plugin in the list requires the name of the plugin and an optional config field specifying proxy level configuration items that will be passed to the plugin upon execution.

For a complete list of available plugins and their corresponding documentation, visit the [documentation for plugins](/docs/v2/plugins).

<div class="example">
<pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
  name: example
spec:
  source:
    path: /foo
  target:
    path: /bar
    backend:
      endpoint: https://foo.bar.com:8443
  plugins:
  - name: apikey
    config:
      bindingName: my-binding
  - name: jwt
    config:
      audienceID: abc123</pre>

</div>

### The `ApiKey` Resource

This resource configures API keys. Note that by itself, an <code>ApiKey</code> resource does not grant permission to any <code>ApiProxy</code>. Permissions are granted via the <code>ApiKeyBinding</code> resource (the next resource we will explore).

> Note that this resource is <i>cluster scoped</i>. This means that resources of this kind are unique per cluster, not per namespace.

Each <code>ApiKey</code> resource specifies a list of revisions. A revision is a specific API key value that may either be active or inactive. The value is both rsa encrypted and base64 encoded. This format caters well for API key rotation.

Below is an example of an <code>ApiKey</code> resource.

<div class="example">
<pre>
---
apiVersion: kanali.io/v2
kind: ApiKey
metadata:
  name: example
spec:
  revisions:
  - data: aGVsbG8=
    status: Active
  - data: d29ybGQ=
    status: Inactive
</div>

### The `ApiKeyBinding` Resource

This resource is responsible for granting course and fine grained permissions to a list of <code>ApiKey</code> resources. Toggle below to learn more about the different options available to you when granting permissions.

<div class="example">
  <nav id="livetabs" data-component="tabs" data-live=".tab-live2"></nav>

  <div id="tab-coarse-grained" data-title="Coarse-grained" class="tab-live2">
  <pre>
---
apiVersion: kanali.io/v2
kind: ApiKeyBinding
metadata:
  name: example
spec:
  keys:
  - name: key-one
    defaultRule:
      global: true
  - name: key-two
    defaultRule:
      granular:
        verbs:
        - GET
  </pre>

  Each API key has a default rule. This can be thought of as coarse grained validation for this API key. This default rule can either grant global access or granular access. Global access will grant an API key access to all http methods whereas granular access allows for only a subset of methods to be granted.
  </div>

  <div id="tab-fine-grained" data-title="Fine-grained" class="tab-live2">
  <pre>
&#8209;&#8209;&#8209;
apiVersion: kanali.io/v2
kind: ApiKeyBinding
metadata:
  name: example
spec:
  keys:
  - name: key-one
    defaultRule:
      global: true
    subpaths:
    - path: /foo
      rule:
        global: false
  </pre>

  Fine-grained access is granted with the <code>subpaths</code> field. This field contains a list of <b>target paths</b> that define a higher priority rule than that of the default rule.

  Let's explore an example using the above <code>ApiKeyBinding</code> resource.

  Suppose a request is made with the <i>key-one</i> API key. Let's also assume that the upstream service will see a request with the path <code>/foo/bar</code>. When Kanali is evaluating which rule is the highest priority, it will do a regular expression match with each path in the subpaths list. If a match is found, that rule is used. If no item in the subpaths list matches, the default rule is used. Hence, in this example, the first item in the subpath list will be a match and hence no access will be granted.
  </div>

  <div id="tab-rate-limiting" data-title="Rate limiting" class="tab-live2">
  <pre>
&#8209;&#8209;&#8209;
apiVersion: kanali.io/v2
kind: ApiKeyBinding
metadata:
  name: example
spec:
  keys:
  - name: key-one
    defaultRule:
      global: true
    rate:
      amount: 5
      unit: second
  </pre>
  <i><b>NOTE:</b> the implementation for rate limiting is not present in</i> <code>v2.0.0</code><i>. It will however be present in an upcoming minor release.</i>

  It you would like to limit the rate of calls a specific API key is allowed to make to a specific upstream service, simplify configure the <code>rate</code> field.
  </div>
</div>

### The `MockTarget` Resource

This resource allows for the configuration of a fake response that will be sent back to the client if the <code>ApiProxy</code> resource that matches a request specifies that a <code>MockTarget</code> should be used as the upstream service.

To configure a <code>MockTarget</code> resource, configure a list of routes. Think of it as if you were developing a an upstream service and each route in your middleware should represent a route item in a <code>MockTarget</code> resource.

<div class="example">
<pre>
&#8209;&#8209;&#8209;
apiVersion: kanali.io/v2
kind: MockTarget
metadata:
  name: example
spec:
  routes:
  - path: /foo
    status: 200
    methods:
    - GET
    - POST
    headers:
      Content-Type: application/json
    body: |-
      {
        "foo": "bar"
      }
</pre>
</div>

