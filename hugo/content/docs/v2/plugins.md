+++
description = "Available Plugins"
title = "Plugins"
date = "2017-04-10T16:43:08+01:00"
draft = false
weight = 100
bref="Available Plugins"
toc = true
+++

### API Key

This plugin performs API key validation for all requests matching a certain <code>ApiProxy</code>. The plugin will ensure that the course and fine grained permissions that are specified in the given <code>ApiKeyBinding</code> resource are enforced. For details on how to configure the <code>ApiKeyBinding</code> resource, read the corresponding documentation [here](/docs/v2/configuration/#the-apikeybinding-resource).

#### Configuration
<table>
  <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  <tr><td><code>bindingName</code></td><td><code>string</code></td><td>Name of <code>ApiKeyBinding</code> resource, in the same namespace as the <code>ApiProxy</code>, to be used.</td></tr>
<table>

#### Example
<pre>
plugins:
- name: apikey
  config:
    bindingName: my-binding</pre>

### JWT

> coming soon