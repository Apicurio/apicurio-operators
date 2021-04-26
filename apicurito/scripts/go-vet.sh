#!/bin/bash

if [[ -z ${CI} ]]; then
    if hash openapi-gen 2>/dev/null; then
        openapi-gen --logtostderr=true -o "" \
            -i ./pkg/apis/apicur/v1alpha1 -O zz_generated.openapi -p ./pkg/apis/apicur/v1alpha1
    else
        echo "skipping go openapi generation"
    fi

    osdk_version=$(operator-sdk version | sed -n 's/.*version: "v\([^"]*\)".*/\1/p')
    if [[ ${osdk_version} == 0.* ]]; then
      echo "operator-sdk >= 1.0.0 required. Please upgrade ..."
      exit 1
    else
      # Calls the config/Makefile which in turn uses controller-gen
      # As described by the operator-sdk documentation
      make -C config generate
    fi
fi
go vet ./...
