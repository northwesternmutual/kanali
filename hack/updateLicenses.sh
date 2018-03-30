#!/bin/bash

set -e

python hack/updateLicense.py $(git ls-files "*\.go")
