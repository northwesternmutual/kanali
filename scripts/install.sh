#!/bin/bash

# silently check if helm is installed
which helm > /dev/null
# install helm if not present
if [ $? != 0 ]; then
   echo "installing helm"
   curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh > /dev/null
   chmod 700 get_helm.sh > /dev/null
   ./get_helm.sh > /dev/null
fi

# deploy tiller
helm init > /dev/null

# add necessary rbac permissions for helm
./scripts/helm-rbac.sh > /dev/null

# add helm repositories for dependencies
helm repo add incubator https://kubernetes-charts-incubator.storage.googleapis.com/ > /dev/null

# wait for the tiller pod to be reader
kubectl rollout status -w deployment/tiller-deploy --namespace=kube-system

# while sleep 1
# do
#     helm install ./helm --name kanali &>/dev/null && break || continue
# done

# install dependencies
helm dep up ./helm > /dev/null

# start kanali and dependencies
helm install ./helm --name kanali > /dev/null

kubectl get pods --all-namespaces