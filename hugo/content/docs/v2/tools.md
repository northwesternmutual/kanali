+++
date = "2017-04-10T16:41:16+01:00"
weight = 70
description = "Developer tooling"
title = "Tools"
draft = false
bref =  "Developer tooling"
toc = true
+++

### Kanalictl

Kanalictl is a command line utility that provides tooling for Kanali.

#### Installation
```
# replace 'darwin' and 'amd64' with your OS and ARCH
$ curl -O https://s3.amazonaws.com/kanalictl/release/$(curl -s https://s3.amazonaws.com/kanalictl/release/latest.txt)/darwin/amd64/kanalictl
$ chmod +x kanalictl
$ sudo mv kanalictl /usr/local/bin/kanalictl
$ kanalictl -h
```

#### Usage
```
$ kanalictl -h
cli interface for kanali

Usage:
  kanalictl [command]

Available Commands:
  apikey      performs operations on API key resources
  version     version
```

#### Commands

##### `apikey`
**Description:** Performs utility operations on <code>ApiKey</code> resources.

<table class="bordered">
  <tr>
    <td><b>subcommand</b></td><td class="w20"><b>description</b></td><td><b>flags</b></td>
  </tr>
  <tr>
    <td><div class="row align-vertical"><code>generate</code></div></td><td>Generates an API key resource.</td><td><table class="unstyled"><tr><td><code>--key.data</code></td><td>Existing API key data</td></tr><tr><td><code>--key.length</code></td><td>Desired length of API key.</td></tr><tr><td><code>--key.name</code></td><td>Name of API key.</td></tr><tr><td><code>--key.out_file</code></td><td>Output file.</td></tr><tr><td><code>--key.public_key_file</code></td><td>Path to RSA public key.</td></tr></table></td>
  </tr>
  <tr>
    <td><code>decrypt</code></td><td>Decrypts one or more API key resources.</td><td><table class="unstyled"><tr><td><code>--key.in_file</code></td><td>Existing API key data</td></tr><tr><td><code>--key.private_key_file</code></td><td>Desired length of API key.</td></tr></table></td>
  </tr>
</table>

##### `version`

**Description:** Displays the current version of Kanalictl.