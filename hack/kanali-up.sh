#!/bin/bash

set -e
source hack/vars.sh

if [ -z "$KANALI_PLUGINS_APIKEY_DECRYPTION_KEY_FILE" ]; then
  echo "creating temporary RSA assets"
  TEMP_DIR=$(mktemp -d -t $PROJECT_NAME)
  mkdir $TEMP_DIR/rsa
  TEMP_RSA_DIR=$TEMP_DIR/rsa
  openssl genrsa -out $TEMP_RSA_DIR/private.pem 2048  > /dev/null 2>&1
  openssl rsa -in $TEMP_RSA_DIR/private.pem -pubout -out $TEMP_RSA_DIR/public.pem  > /dev/null 2>&1
  echo "temporary RSA assets created at ${TEMP_DIR}/rsa"
  export KANALI_PLUGINS_APIKEY_DECRYPTION_KEY_FILE=$TEMP_RSA_DIR/private.pem

  # If OSX and using ZSH as shell, export environment variable there
  if [ -f $HOME/.zshrc ] && ! $(cat ~/.zshrc | grep KANALI_PLUGINS_APIKEY_DECRYPTION_KEY_FILE); then
    echo "export KANALI_PLUGINS_APIKEY_DECRYPTION_KEY_FILE=${TEMP_RSA_DIR}/private.pem" >> $HOME/.zshrc
  fi
fi

kanali start \
  --process.log_level="debug" \
  --server.insecure_port=8080 \
  --server.insecure_bind_address="0.0.0.0" \
  --kubernetes.kubeconfig=$HOME/.kube/config \
  --profiling.enabled \
  --profiling.insecure_port=9090 \
  --profiling.insecure_bind_address="0.0.0.0" \
  --prometheus.insecure_port=9000 \
  --prometheus.insecure_bind_address="0.0.0.0" \
  --proxy.enable_cluster_ip=true \
  --proxy.enable_mock_responses=true \
  --proxy.tls_common_name_validation=true \
  --plugins.apiKey.header_key="apikey" \
  --plugins.apiKey.decryption_key_file=$KANALI_PLUGINS_APIKEY_DECRYPTION_KEY_FILE
