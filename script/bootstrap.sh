#!/usr/bin/env sh
#
# bootstrap.sh will check for and install any dependencies we have for building and using iofogctl
#
# Usage: ./bootstrap.sh
#

set -e

#
# All our Go related stuff
#

# Is go installed?
if ! checkForInstallation "go"; then
    echoNotify "\nYou do not have Go installed. Please install and re-run bootstrap."
    exit 1
fi

# Is dep installed?
if ! checkForInstallation "dep"; then
    echoInfo " Attempting to install 'go dep'"
    go get -u github.com/golang/dep/cmd/dep
fi
