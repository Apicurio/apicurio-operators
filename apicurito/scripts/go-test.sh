#!/bin/bash

export GO111MODULE=on
if [[ -z ${CI} ]]; then
    ./scripts/go-vet.sh
    ./scripts/go-fmt.sh
fi

GOFLAGS="" go test -test.short -mod=vendor ./cmd/... ./pkg/...