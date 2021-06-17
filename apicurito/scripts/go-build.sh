#!/bin/bash
REGISTRY=quay.io/${USER}
IMAGE=apicurito-operator
TAG=v1.1.0
GIT_COMMIT=$(git rev-list -1 HEAD)
export GO111MODULE=on
GOFLAGS="-X github.com/apicurio/apicurio-operators/apicurito/pkg/cmd.GitCommit=${GIT_COMMIT}"

go mod vendor
go generate ./...

./scripts/go-test.sh

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags "${GOFLAGS}" \
  -o build/_output/bin/apicurito \
  -mod=vendor github.com/apicurio/apicurio-operators/apicurito/cmd/manager

docker build . -t $REGISTRY/$IMAGE:$TAG
