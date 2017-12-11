#!/bin/bash

# silently check if helm is installed
which helm > /dev/null
# install helm if not present
if [ $? != 0 ]; then
   echo "installing helm"
   curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh
   chmod 700 get_helm.sh
   ./get_helm.sh
fi

# deploy tiller
helm init > /dev/null

# add necessary rbac permissions for helm
./scripts/helm-rbac.sh > /dev/null

# add helm repositories for dependencies
helm repo add incubator https://kubernetes-charts-incubator.storage.googleapis.com/ > /dev/null

# while sleep 1
# do
#     helm install ./helm --name kanali &>/dev/null && break || continue
# done

helm install ./helm --name kanali

kubectl get pods --all-namespaces