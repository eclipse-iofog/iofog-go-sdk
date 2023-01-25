#!/bin/bash

# Get variables
. test/conf/vars.bash

# Delete namespace
iofogctl delete namespace "$NS" --force -v
iofogctl disconnect -n "$NS"
