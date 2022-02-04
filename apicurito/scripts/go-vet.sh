#!/bin/bash

OPENAPI_GEN_VERSION="release-1.22"

if [[ -z ${CI} ]]; then
    if ! command -v openapi-gen &> /dev/null
    then
      echo "Downloading and installing openapi-gen ..."
      go install k8s.io/kube-openapi/cmd/openapi-gen@release-1.22
      if [ $? != 0 ]; then
        echo "Error: Failed to install openapi-gen"
        exit 1
      fi
    else
      echo "Warning: openapi-gen already installed but cannot guarantee if its version is compatible."
      echo "         To ensure compatibility, please set aside the current version."
      exit 1
    fi

    openapi-gen --logtostderr=true -o "" \
        -i ./pkg/apis/apicur/v1alpha1 -O zz_generated.openapi -p ./pkg/apis/apicur/v1alpha1
    if [ $? != 0 ]; then
        echo "Error: openapi-gen failed to generate the API"
        exit 1
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
