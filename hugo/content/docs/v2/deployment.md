+++
date = "2017-04-10T16:41:54+01:00"
weight = 110
description = "Learn how to deploy Kanali"
title = "Deployment"
draft = false
bref= "Learn how to deploy Kanali"
toc = true
+++

### Helm

> This is the recommended deployment option.

[Helm](https://helm.sh) is a package manager for Kubernetes. Helm makes it easy to deploy all the components of an application and its dependencies.

Kanali's Helm chart can be found [here](https://github.com/northwesternmutual/kanali/tree/master/helm/kanali). Before demonstrating how to use this chart, let's bootstrap a local Kubernetes environment so that you can follow along.

```
# Start a local Kubernetes cluster.
$ minikube start --kubernetes-version=v1.9.4 --feature-gates CustomResourceValidation=true

# Update a potentially stale kubeconfig context.
$ minikube update-context
Kubeconfig IP correctly configured, pointing at 192.168.64.24

# Verify that our cluster is ready.
$ kubectl get nodes
NAME       STATUS    ROLES     AGE       VERSION
minikube   Ready     <none>    17h       v1.9.4
```

Now that we have have a running cluster, the Kanali Helm chart can be deployed by executing the following commands.

```
# Retrieve a local copy of Kanali. This is necessary because Kanali does not
# yet live in a chart repository.
$ git clone https://github.com/northwesternmutual/kanali.git

# To ensure that Helm has the appropriate permissions to deploy our chart,
# we need to create an RBAC policy and a service account.
$ kubectl apply -f kanali/hack/helm-rbac.yaml

# Bootstrap helm using the service account created in the previous step.
$ helm init --service-account tiller

# One of the optional dependencies for Kanali is Jaeger, whose chart lives in the following repo.
$ helm repo add incubator https://kubernetes-charts-incubator.storage.googleapis.com/

# Install any required dependencies.
$ helm dep up kanali/helm/kanali

# Install the Kanali chart (note the relative path to the chart).
$ helm install --name kanali kanali/helm/kanali
```

Kanali is now deployed. Note that the default set of Helm values were used. You will probably want to overwrite some of these values for your deployment. All of the values that you can customize can be found [here](https://github.com/northwesternmutual/kanali/tree/master/helm/kanali). Instructions on how to set your custom values can be found [here](https://docs.helm.sh/developing_charts/#charts).

We can wait for it to complete its bootstrapping process with the following command.

```
$ kubectl rollout status -w deployment/kanali --namespace=default
```

Now that Kanali is ready to be used, we can test our installation with the following command. By default, Kanali will bootstrap itself with a set of self signed certificates which is why our https request is insecure.

```
$ curl --insecure $(minikube service kanali-gateway --url --format="https://{{.IP}}:{{.Port}}")
{
  "status": 404,
  "message": "No ApiProxy resource was not found that matches the request.",
  "code": 0,
  "details": "Visit https://northwesternmutual.github.io/kanali/docs/v2/errorcodes for more details."
}
```

**Configurations!** You have successfully bootstrapped Kanali using Helm. To learn more, visit the [tutorial](https://northwesternmutual.github.io/kanali/tutorial).

### Manual

If you like to deploy Kanali a different way, there are just a few things you need to know.

Kanali's artifact is a [Docker](https://www.docker.com/) image. When starting this container, you will need to specify your own configuration arguments if you would like to overwrite the defaults. This can be accomplished by specifying one or more flags, environment variables and/or configuration file updates. Specific details for these configuration items can be found [here](https://northwesternmutual.github.io/kanali/docs/v2/flags).
