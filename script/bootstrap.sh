#!/usr/bin/env sh
#
# bootstrap.sh will check for and install any dependencies we have for building / developing iofog-go-sdk
#
# Usage: ./bootstrap.sh
#

set -e

# Is go installed?
if [ -z $(command -v go) ]; then
    echo "\nYou do not have Go installed. Please install and re-run bootstrap."
    exit 1
fi

# Is go lint installed?
if [ ! "$(command -v golangci-lint)" ]; then
    if [ "$(uname -s)" = "Darwin" ]; then
      brew install golangci-lint
      brew upgrade golangci-lint
    else
      go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
    fi
fi

# Leave gengo last!!

# Is gengo installed?
if [ -z $(command -v deepcopy-gen) ]; then
    echo " Attempting to install 'gengo'"
    go install -v k8s.io/code-generator/cmd/deepcopy-gen@v0.26
#    go install -mod=vendor k8s.io/gengo/examples/deepcopy-gen/
fi

