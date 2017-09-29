#!/bin/bash

LIGHT_BLUE='\033[1;34m'
LIGHT_GREEN='\033[1;32m'
NC='\033[0m'

echo -e "${LIGHT_BLUE}Step 1: Verify that Helm is installed${NC}"
which helm > /dev/null || $(curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh | bash)

echo -e "${LIGHT_BLUE}Step 2: Deploy Helm"
helm init > /dev/null

echo -e "${LIGHT_BLUE}Step 3: Patch RBAC for Helm${NC}"
./scripts/helm-permissions.sh > /dev/null

echo -e "${LIGHT_BLUE}Step 4: Install Kanali, Grafana, InfluxDB, and Jaeger (may take a few minutes)${NC}"

while sleep 1
do
    helm install ./helm --name kanali &>/dev/null && break || continue
done

echo -e "${LIGHT_GREEN}Deployment Complete!${NC}"
