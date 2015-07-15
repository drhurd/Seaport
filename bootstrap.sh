#!/usr/bin/env bash

#export DEBIAN_FRONTEND=noninteractive
sudo apt-get update

# download/install go
mkdir /home/vagrant/golang
wget -q -P /home/vagrant/golang/ https://storage.googleapis.com/golang/go1.4.2.linux-386.tar.gz
tar -xzf /home/vagrant/golang/go1.4.2.linux-386.tar.gz -C /home/vagrant/golang/
echo "export GOROOT=/home/vagrant/golang/go" >> "/home/vagrant/.profile"
echo "export PATH=\$PATH:\$GOROOT/bin" >> "/home/vagrant/.profile"
