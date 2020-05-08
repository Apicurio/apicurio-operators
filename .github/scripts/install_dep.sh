#!/bin/bash
set -e

echo "Update archives"

sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt-get update

echo "Install golang"
sudo apt-get install golang-1.13
export GO111MODULE='on'

echo "Install operator sdk"
curl -fsSL https://github.com/operator-framework/operator-sdk/releases/download/v0.17.0/operator-sdk-v0.17.0-x86_64-linux-gnu > operator-sdk
chmod +x operator-sdk
sudo mv operator-sdk /usr/local/bin/operator-sdk

echo "Clean cache"
sudo apt-get clean

#Docker installation from google default repositories
sudo apt-get -y update

sudo apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common
sudo mkdir -p /mnt/docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"

sudo apt-get -y update
sudo apt-get install docker-ce docker-ce-cli containerd.io