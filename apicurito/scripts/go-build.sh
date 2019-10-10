#!/bin/sh
REGISTRY=quay.io/lgarciaac
IMAGE=apicurito-operator
TAG=v0.1
CFLAGS="--redhat --build-tech-preview"

go generate ./...
if [[ -z ${CI} ]]; then
    ./scripts/go-test.sh
    operator-sdk build ${REGISTRY}/${IMAGE}:${TAG}
   
else
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o build/_output/bin/apicurito github.com/apicurio/apicurio-operators/apicurito/cmd/manager

fi



