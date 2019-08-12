#!/bin/bash
#
# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# This script is intended to install a docker-based gateway deployment

set -e

CWAG="cwag"
FEG="feg"
INSTALL_DIR="/tmp/magmagw_install"

# TODO: Update docker-compose to stable version

# Using RC as opposed to stable (1.24.0) due to
# SCTP port mapping support
DOCKER_COMPOSE_VERSION=1.25.0-rc1

DIR="."
echo "Setting working directory as: $DIR"
cd "$DIR"

if [ -z $1 ]; then
  echo "Please supply a gateway type to install. Valid types are: ['$FEG', '$CWAG']"
  exit
fi

GW_TYPE=$1
echo "Setting gateway type as: '$GW_TYPE'"

if [ "$GW_TYPE" != "$FEG" ] && [ "$GW_TYPE" != "$CWAG" ]; then
  echo "Gateway type '$GW_TYPE' is not valid. Valid types are: ['$FEG', '$CWAG']"
  exit
fi

# Ensure necessary files are in place
if [ ! -f /opt/magma/env/.env ]; then
    echo ".env file is missing! Please add this file to the directory that you are running this command and re-try."
    exit
fi

if [ ! -f /opt/magma/certs/rootCA.pem ]; then
    echo "rootCA.pem file is missing! Please add this file to the directory that you are running this command and re-try."
    exit
fi

# TODO: Remove this once .env is used for control_proxy
if [ ! -f /opt/magma/env/control_proxy.yml ]; then
    echo "control_proxy.yml file is missing! Please add this file to the directory that you are running this command and re-try."
    exit
fi

# Fetch files from github repo
rm -rf "$INSTALL_DIR"
mkdir -p "$INSTALL_DIR"

MAGMA_GITHUB_URL="{{ .Values.cwf.repo.url }}"
git -C "$INSTALL_DIR" clone "$MAGMA_GITHUB_URL" -b {{ .Values.cwf.repo.branch }}

# TODO: Add this back once this code is included in a github version
#TAG=$(git -C $INSTALL_DIR/magma tag | tail -1)
#git -C $INSTALL_DIR/magma checkout "tags/$TAG"

if [ "$GW_TYPE" == "$CWAG" ]; then
  MODULE_DIR="cwf"

  # Run CWAG ansible role to setup OVS
  echo "Copying and running ansible..."
  apt-add-repository -y ppa:ansible/ansible
  apt-get update -y
  apt-get -y install ansible
  ANSIBLE_CONFIG="$INSTALL_DIR"/magma/"$MODULE_DIR"/gateway/ansible.cfg ansible-playbook "$INSTALL_DIR"/magma/"$MODULE_DIR"/gateway/deploy/cwag.yml -i "localhost," -c local -v
fi

if [ "$GW_TYPE" == "$FEG" ]; then
  MODULE_DIR="$GW_TYPE"

  # Load kernel module necessary for docker SCTP support
  sudo modprobe nf_conntrack_proto_sctp
fi

cp "$INSTALL_DIR"/magma/"$MODULE_DIR"/gateway/docker/docker-compose.yml .
cp "$INSTALL_DIR"/magma/orc8r/tools/docker/upgrade_gateway.sh .
# Install Docker
sudo apt-get update
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

# Install Docker-Compose
sudo curl -L "https://github.com/docker/compose/releases/download/"$DOCKER_COMPOSE_VERSION"/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create snowflake to be mounted into containers
touch /etc/snowflake

echo "Placing configs in the appropriate place..."
mkdir -p /var/opt/magma
mkdir -p /var/opt/magma/configs
mkdir -p /var/opt/magma/certs
mkdir -p /etc/magma
mkdir -p /var/opt/magma/docker

# Copy default configs directory
cp -TR "$INSTALL_DIR"/magma/"$MODULE_DIR"/gateway/configs /etc/magma

# Copy config templates
cp -R "$INSTALL_DIR"/magma/orc8r/gateway/configs/templates /etc/magma

# Copy certs
cp /opt/magma/certs/rootCA.pem /var/opt/magma/certs/

# Copy control_proxy override
cp /opt/magma/env/control_proxy.yml /etc/magma/

# Copy docker files
cp docker-compose.yml /var/opt/magma/docker/
cp /opt/magma/env/.env /var/opt/magma/docker/

# Copy upgrade script for future usage
cp upgrade_gateway.sh /var/opt/magma/docker/

cd /var/opt/magma/docker
source .env

{{- if and .Values.cwf.image.username .Values.cwf.image.password }}
echo "Logging into docker registry at $DOCKER_REGISTRY"
docker login "$DOCKER_REGISTRY" --username {{ .Values.cwf.image.username }} --password {{ .Values.cwf.image.password }}
{{- end }}
docker-compose pull
docker-compose -f docker-compose.yml up -d

echo "Installed successfully!!"
