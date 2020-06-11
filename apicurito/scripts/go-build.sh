#!/bin/bash
REGISTRY=quay.io/${USER}
IMAGE=apicurito-operator
TAG=v0.2

export GO111MODULE=on
go mod vendor

go generate ./...
if [[ -z ${CI} ]]; then
    ./scripts/go-test.sh
    operator-sdk build ${REGISTRY}/${IMAGE}:${TAG}
   
else
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o build/_output/bin/apicurito -mod=vendor github.com/apicurio/apicurio-operators/apicurito/cmd/manager

fi