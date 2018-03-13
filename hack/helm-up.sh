#!/bin/bash

set -e

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

echo "Helm deployment of Kanali successful!"