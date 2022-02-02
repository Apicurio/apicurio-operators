#!/bin/bash

if [ -z "${IMAGE}" ]; then
  echo "Error: IMAGE env var not defined"
  exit 1
fi

if [ -z "${TAG}" ]; then
  echo "Error: TAG env var not defined"
  exit 1
fi

export GO111MODULE=on
go mod vendor

./scripts/go-test.sh

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -a \
  -o build/_output/bin/apicurito \
  -mod=vendor github.com/apicurio/apicurio-operators/apicurito/cmd/manager
if [ $? != 0 ]; then
  echo "Error: build failed"
  exit 1
fi

echo
echo "=== Building image ..."
docker build . -t ${IMAGE}:${TAG}
