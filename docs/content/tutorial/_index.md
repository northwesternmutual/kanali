+++
description = "Add some motion, shaking, pulsing, sliding and more"
title = "Quick Start"
date = "2017-04-10T16:43:08+01:00"
draft = false
weight = 200
bref="foo"
type = "tutorial"
+++

<div class="example">

<nav id="tabs" class="tabs" data-component="tabs">
    <ul class="hide">
        <li class="hidden active"><a href="#tab1">one</a></li>
        <li><a href="#tab2">two</a></li>
        <li><a href="#tab3">three</a></li>
        <li><a href="#tab4">four</a></li>
        <li><a href="#tab5">five</a></li>
    </ul>
</nav>

<div id="tab1">
<p>
foo bar
</p>
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
<div id="tab2">second</div>
<div id="tab3">third</div>
<div id="tab4">fourth</div>
<div id="tab5">fifth</div>

<br />


<div class="group">
<button class="float-left button outline big" onclick="$('#tabs').tabs('prev');">prev</button>
<button class="float-right button outline big" onclick="$('#tabs').tabs('next');">next</button>
</div>

</div>