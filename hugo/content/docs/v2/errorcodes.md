+++
date = "2017-04-10T16:41:54+01:00"
weight = 30
description = "Details around Kanali's error codes"
title = "Error Codes"
draft = false
bref= "Details around Kanali's error codes"
toc = true
+++

### Introduction
If you're visiting this page, changes are that Kanali has given you a response similar to this one:

```json
{
  "status": 404,
  "message": "No ApiProxy resource was not found that matches the request.",
  "code": 0,
  "details": "Visit https://northwesternmutual.github.io/kanali/docs/v2/errorcodes/#00 for more details."
}
```
As promised in that response, the table below contains detailed information about each unique error code.

### Details

#### `00`

**Status:** <span class="label error">404</span>

**Message:** *No ApiProxy resource was not found that matches the request.*

**Description:** Kanali was unable to match the given request to an <code>ApiProxy</code> resource at the time the request occurred.

**Troubleshooting:** Verify the presence of the <code>ApiProxy</code> resource inside of the Kubernetes cluster that Kanali is deployed to. In addition, verify that the path in the request is prefixed with the source path of the <code>ApiProxy</code> resource.

#### `01`

**Status:** <span class="label error">500</span>

**Message:** *An unknown error occurred.*

**Description:** Something went wrong. However, no details are known at this time.

**Troubleshooting:** Reach out to your Kanali administrator for further troubleshooting.

#### `02`

**Status:** <span class="label error">404</span>

**Message:** *No MockTarget resource was not found that matches the request.*

**Description:** Kanali was unable to find the <code>MockTarget</code> resource that was specified in the <code>ApiProxy</code> resource that matched this request.

**Troubleshooting:** Verify that the name of the <code>MockTarget</code> that is configured in matching <code>ApiProxy</code> resource exists in the same namespace as the <code>ApiProxy</code> resource.

#### `03`

**Status:** <span class="label error">500</span>

**Message:** *Could not open or load plugin.*

**Description:** Kanali was unable to open or load one of the plugins specified in the <code>ApiProxy</code> resource that matched this request.

**Troubleshooting:** Verify that there is a <code>.so</code> file for each corresponding plugin used by the matching <code>ApiProxy</code>. Look for these files in the location specified by the <code>--plugins.location</code> configuration flag. This value can be found in one of the very first logs that Kanali logs upon startup.

#### `04`

**Status:** <span class="label error">500</span>

**Message:** *Could not lookup plugin symbol.*

**Description:** Kanali attempted to execute one of the plugins configured on the <code>ApiProxy</code> resource that matched this request. However, when attempting to load one of the plugins, no exported variable, <code>Plugin</code> was found.

**Troubleshooting:** Confirm that all of the plugins that the matching <code>ApiProxy</code> resource is using export the variable <code>Plugin</code>. If you find that one of the plugins does not do this, reach out to your Kanali administrator.

#### `05`

**Status:** <span class="label error">500</span>

**Message:** *Plugin does not implement the correct interface.*

**Description:** Kanali attempted to execute one of the plugins configured on the <code>ApiProxy</code> resource that matched this request. However, when attempting to execute the plugin, it was determined that the plugin did not implement the required interface. This interface can be found [here](https://godoc.org/github.com/northwesternmutual/kanali/pkg/plugin#Plugin).

**Troubleshooting:** Confirm that all of the plugins that the matching <code>ApiProxy</code> resource is using implement the interface described above. If you find that one of the plugins do not implement this interface, reach out to your Kanali administrator.

#### `06`

**Status:** <span class="label error">500</span>

**Message:** *Could not retrieve Kubernetes secret.*

**Description:** Kanali attempted to retrieve the secret specified in the <code>ssl</code> field of the <code>ApiProxy</code> resource that matched this request. Due to one of a few reasons, there was an issue retrieving and/or using this secret.

**Troubleshooting:** First, verify that there is a secret that has the same name as the one specified in the <code>ApiProxy</code> resource that matched this request. Note that this secret must also live in the same namespace as the <code>ApiProxy</code> resource. If the secret exists, verify that the <code>kanali.io/enabled: 'true'</code> annotation is present on the secret. Next, verify that the secret contains the correct data fields. Documentation about what fields are required can be found [here](/docs/v2/configuration/#ssl). If everything checks out, contact your Kanali administrator and let them know that Kanali may be having trouble connecting to the Kubernetes API server.

#### `07`

**Status:** <span class="label error">500</span>

**Message:** *Could not create x509 key pair.*

**Description:** Kanali found a valid secret as configured in the <code>ApiProxy</code> resource that matched this request. However, while attempting to parse the public/private key pair from the data in this secret, an error occurred.

**Troubleshooting:** Validate that the public/private key pair that is specified in the above described secret is valid. One way to verify this is to use the following commands to generate <code>md5</code> hashes. If the hashes are different, your key pair is not valid.

<pre>
openssl x509 -noout -modulus -in server.crt | openssl md5
openssl rsa -noout -modulus -in server.key | openssl md5
</pre>

#### `08`

**Status:** <span class="label error">502</span>

**Message:** *Could not get a valid or any response from the upstream server.*

**Description:** Kanali attempted to issue is request to an upstream service. However, something with this request went wrong. The issue could stem from the configuration of the <code>ApiProxy</code> resource matching this request, from the upstream service, from the networking in between Kanali and the upstream, or something else.

**Troubleshooting:** First, verify that you are able to access the upstream service without Kanali (e.g. using curl from a busybox pod). If this is successful, make sure that the <code>ApiProxy</code> resource mathing this request is correctly configured so that it will route traffic to the expected upstream service. You might also ensure that any networking related items such as route tables, firewalls, etc. are correctly configured.

#### `09`

**Status:** <span class="label error">500</span>

**Message:** *Could not retrieve Kubernetes services.*

**Description:** Kanali determined that for this request, the upstream service is a Kubernetes service. However, when attempting to find out which service it might be, an error occurred. This does not imply that Kanali is unable to communicate with the Kubernetes API server. However, it does imply that Kanali has an issue.

**Troubleshooting:** Try recycling Kanali's pod(s) and reach out to your Kanali administrator.

#### `10`

**Status:** <span class="label error">500</span>

**Message:** *Could not retrieve any matching Kubernetes services.*

**Description:** Kanali determined that for this request, the upstream service is a Kubernetes service. However, when attempting to find out which service it might be, none were found.

**Troubleshooting:** Review the documentation for configuring the <code>ApiProxy</code> resource [here](/docs/v2/configuration). If static service discovery is being used, verify that a Kubernetes service in the same namespace as the matching <code>ApiProxy</code> resource exists. If dynamic service discovery is being used, verify that there exists a service in that same namespace that contains the dynamic labels.

#### `11`

**Status:** <span class="label error">500</span>

**Message:** *Plugin threw a runtime error.*

**Description:** Kanali attempted to execute one of the plugins configured on the <code>ApiProxy</code> resource that matched this request. However, when attempting to execute one of the plugins, a runtime error was thrown causing the plugin to crash.

**Troubleshooting:** Contact your Kanali administrator as a bug in a plugin has been identified.

#### `12`

**Status:** <span class="label error">403</span>

**Message:** *My lips are sealed.*

**Description:** All I can say is that something went wrong.

**Troubleshooting:** Like the message says, *My lips are sealed.*...

#### `13`

**Status:** <span class="label error">401</span>

**Message:** *Api key is not authorized.*

**Description:** Kanali was able to associate this request to a matching <code>ApiProxy</code> resource. In addition, the API key plugin was able to find an <code>ApiKeyBinding</code> resource that was configured in the <code>ApiProxy</code> resource. However, in the <code>ApiKeyBinding</code> resource, access it not granted to the API key that was used for the type of request performed.

**Troubleshooting:** Verify that the API key used does not have permission to perform the request. If the <code>ApiKeyBinding</code> resource grants the API key used access to the specific request performed, contact your Kanali administrator for further troubleshooting.

#### `14`

**Status:** <span class="label error">429</span>

**Message:** *The Api key you are using has exceeded its rate limit.*

**Description:** Kanali was able to associate this request to a matching <code>ApiProxy</code> resource. In addition, the API key plugin was able to find an <code>ApiKeyBinding</code> resource that was configured in the <code>ApiProxy</code> resource. However, the API key used has exceeded its rate limit or quota policy.

**Troubleshooting:** Verify that the API key used did indeed exceed its rate limit or quota. If the rate limit or quota was not exceeded, contact your Kanali administrator for further troubleshooting.