#!/bin/bash
REGISTRY=quay.io/${USER}
IMAGE=apicurito-operator
TAG=v1.1.0
GIT_COMMIT=$(git rev-list -1 HEAD)
FUSE_VERSION=$(grep ^DEFAULT_VERSION ./config/vars/Makefile|awk '{print $3}')
FUSE_PREV_VERSION=$(grep ^DEFAULT_PREVIOUS_VERSION ./config/vars/Makefile|awk '{print $3}')
export GO111MODULE=on
GOFLAGS="-X github.com/apicurio/apicurio-operators/apicurito/pkg/cmd.GitCommit=${GIT_COMMIT}"
GOFLAGS=$GOFLAGS" -X github.com/apicurio/apicurio-operators/apicurito/version.Version=${FUSE_VERSION}"
GOFLAGS=$GOFLAGS" -X github.com/apicurio/apicurio-operators/apicurito/version.PriorVersion=${FUSE_PREV_VERSION}"

go mod vendor
go generate ./...

./scripts/go-test.sh

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags "${GOFLAGS}" \
  -o build/_output/bin/apicurito \
  -mod=vendor github.com/apicurio/apicurio-operators/apicurito/cmd/manager

docker build . -t $REGISTRY/$IMAGE:$TAG
