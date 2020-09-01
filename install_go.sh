#!/bin/bash

curl -O https://dl.google.com/go/go1.15.linux-amd64.tar.gz

sha256sum go1.15.linux-amd64.tar.gz

tar xvf go1.15.linux-amd64.tar.gz

chown -R root:root ./go

mv go /usr/local

ln -s /usr/local/go /usr/local/bin/go
ln -s /usr/local/gofmt /usr/local/bin/gofmt

GOPATH=$HOME/goapp
PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

go version

