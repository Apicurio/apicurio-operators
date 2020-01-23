#!/bin/bash

if [[ -z ${CI} ]]; then
    operator-sdk generate k8s
    operator-sdk generate openapi
fi
go vet ./...