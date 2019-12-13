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

# Leave gengo last!!

# Is gengo installed?
if [ -z $(command -v deepcopy-gen) ]; then
    echo " Attempting to install 'gengo'"
    go install -mod=vendor k8s.io/gengo/examples/deepcopy-gen/
fi

