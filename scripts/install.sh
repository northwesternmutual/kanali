#!/bin/bash

# check if helm is installed
which helm

# install helm if not present
if [ $? != 0 ]; then
  set -e
  curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh
  chmod 700 get_helm.sh
  ./get_helm.sh
fi

set -e

# deploy tiller
helm init

# add necessary rbac permissions for helm
./scripts/helm-rbac.sh

# add helm repositories for dependencies
helm repo add incubator https://kubernetes-charts-incubator.storage.googleapis.com/

# wait for the tiller deployment to be ready
kubectl rollout status -w deployment/tiller-deploy --namespace=kube-system

# install dependencies
helm dep up ./helm

# start kanali and dependencies
helm install ./helm --name kanali --set kanali.tag=${COMMIT:-latest}

# wait for deployments to be ready
kubectl rollout status -w deployment/kube-dns --namespace=kube-system
# kubectl rollout status -w deployment/etcd --namespace=default
kubectl rollout status -w deployment/kanali --namespace=default

kubectl get pods --all-namespaces