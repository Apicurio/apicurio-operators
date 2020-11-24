#!/bin/bash
REGISTRY=quay.io/${USER}
IMAGE=apicurito-operator
TAG=v1.1.0

export GO111MODULE=on
go mod vendor

go generate ./...
./scripts/go-test.sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o build/_output/bin/apicurito -mod=vendor github.com/apicurio/apicurio-operators/apicurito/cmd/manager
docker build . -t $REGISTRY/$IMAGE:$TAG