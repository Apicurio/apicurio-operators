#!/bin/bash

if [[ -z ${CI} ]]; then
    if hash openapi-gen 2>/dev/null; then
        openapi-gen --logtostderr=true -o "" \
            -i ./pkg/apis/apicur/v1alpha1 -O zz_generated.openapi -p ./pkg/apis/apicur/v1alpha1
    else
        echo "skipping go openapi generation"
    fi
fi
go vet ./...