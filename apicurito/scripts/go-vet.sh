#!/bin/sh

if [[ -z ${CI} ]]; then
    ./scripts/go-dep.sh
    operator-sdk generate k8s
fi
go vet ./...