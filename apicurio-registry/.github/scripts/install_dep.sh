#!/bin/bash

echo "Update archives"

sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt-get update

echo "Install golang"
sudo apt-get install golang-1.13

echo "Clean cache"
sudo apt-get clean

set -e

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