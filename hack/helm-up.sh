#!/bin/bash

set -e

# Ensure that dependent binaries are present
which helm || (echo "Helm not found in your path. Please install before continuing." && exit 1)
which kubectl || (echo "Kubectl not found in your path. Please install before continuing." && exit 1)
which minikube || (echo "Minikube not found in your path. Please install before continuing." && exit 1)

echo "Bootstraping a local Kubernetes environment...."
minikube start

echo "Adding the Helm incubator repo as some of our optional dependencies live here....."
helm repo add incubator https://kubernetes-charts-incubator.storage.googleapis.com/
echo "Applying RBAC policies for Helm....."
kubectl apply -f ./hack/helm-rbac.yaml

echo "Bootstrapping Helm....."
helm init --service-account tiller

echo "Waiting for Helm to become ready....."
kubectl rollout status -w deployment/tiller-deploy --namespace=kube-system

echo "Installing our optional dependencies....."
helm dep up ./helm

echo "Deploying Kanali....."
helm install ./helm --name kanali

echo "Waiting for Kanali to become ready....."
kubectl rollout status -w deployment/kanali --namespace=default

echo "Helm deployment of Kanali successful!"