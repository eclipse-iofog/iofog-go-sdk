#!/bin/bash

set -e

# iofogctl
curl -s https://packagecloud.io/install/repositories/iofog/iofogctl/script.deb.sh | sudo bash
sudo apt install -qy iofogctl
iofogctl version
