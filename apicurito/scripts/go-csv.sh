#!/bin/sh

if [ -z "$CONFIG_FILE" ]
then
      echo "\$CONFIG_FILE is empty"
      exit 1
fi

go run -ldflags="-X github.com/apicurio/apicurio-operators/apicurito/pkg/configuration.ConfigFile=$CONFIG_FILE" ./tools/csv-gen/main.go