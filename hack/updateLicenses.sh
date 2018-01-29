#!/bin/bash

set -e

python hack/updateLicense.py $(git ls-files "*\.go" | grep -v thrift-gen | grep -v tracetest)
