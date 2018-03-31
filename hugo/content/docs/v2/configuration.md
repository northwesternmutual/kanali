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
The goal of this section is to provide an introduction to each of Kanali's configurable resources. More in depth documentation for each resource can be found <a>here</a>.

<br/>

### The `ApiProxy` Resource

This resource declaratively defines how your upstream services are 

<div class="example">
  <nav id="livetabs" data-component="tabs" data-live=".tab-live"></nav>

  <div id="tab-service-static" data-title="Service (static)" class="tab-live">
    <pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
 name: example
 namespace: default
spec:
 source:
   path: /example
 target:
   service:
     name: serviceName
     port: 8080
    </pre>
  </div>
  <div id="tab-service-dynamic" data-title="Service (dynamic)" class="tab-live">
    <pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
 name: example
 namespace: default
spec:
 source:
   path: /example
 target:
   service:
     port: 8080
     labels:
     - name: key
       value: value
     - name: deploy
       header: x-foo-deployment
    </pre>
  </div>
  <div id="tab-endpoint" data-title="Endpoint" class="tab-live">
  <pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
 name: example
 namespace: default
spec:
 source:
   path: /example
 target:
   backend:
     endpoint: https://foo.bar.com:8443
  </pre>
</div>
  <div id="tab-mock" data-title="Mock" class="tab-live">
  <pre>
---
apiVersion: kanali.io/v2
kind: ApiProxy
metadata:
 name: example
 namespace: default
spec:
 source:
   path: /example
 target:
   backend:
     mock:
       mockTargetName: mockTargetName
  </pre>
</div>
</div>

### The `ApiKey` Resource

### The `ApiKeyBinding` Resource

### The `MockTarget` Resource

