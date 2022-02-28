CONTROLLER_GEN_VERSION := v0.4.1
OPERATOR_SDK_VERSION := v1.5.0
KUSTOMIZE_VERSION := v4.1.2

PACKAGE := fuse-apicurito

DEFAULT_VERSION := 7.10.1
DEFAULT_PREVIOUS_VERSION := 7.10.0
DEFAULT_TAG := 1.10
DEFAULT_IMAGE := registry-proxy.engineering.redhat.com/rh-osbs/fuse7-fuse-apicurito-rhel8-operator
DEFAULT_NAMESPACE := $(shell $(CONFIG_DIR)/script/namespace.sh)
DEFAULT_CHANNEL := $(PACKAGE)-$(shell v=$(DEFAULT_VERSION) && echo "$${v%\.[0-9]}.x")
DEFAULT_CSV := manifests/bases/apicurito.clusterserviceversion.yaml
BUNDLE_INFIX := manifest

CSV_SUPPORT := Fuse Apicurito
CSV_DISPLAY_NAME := Red Hat Integration - API Designer

#
# Vars that can be overridden by external env vars
#
VERSION ?= $(DEFAULT_VERSION)
PREVIOUS_VERSION ?= $(DEFAULT_PREVIOUS_VERSION)
DEBUG ?= false
KUBE_USER ?= developer
DEFAULT_IMAGE ?=
IMAGE ?= $(DEFAULT_IMAGE)
DEFAULT_TAG ?= 1.10
TAG ?= $(DEFAULT_TAG)
LEGACY ?= false
NAMESPACE ?= $(DEFAULT_NAMESPACE)

CHANNELS ?= $(DEFAULT_CHANNEL)
# Default bundle image tag
BUNDLE_IMG ?= $(IMAGE)-$(BUNDLE_INFIX):$(VERSION)
# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
ifneq ($(origin PACKAGE), undefined)
BUNDLE_PACKAGE := --package=$(PACKAGE)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL) $(BUNDLE_PACKAGE)

BUNDLE_DIR := bundle
PKG := ../pkg
CRD := crd/bases
MANAGER := manager
MANIFESTS := manifests
RBAC := rbac
TIMESTAMP=$(shell date -u '+%F %T %Z')
DOMAIN := apicur.io

CSV_REPLACES := $(PACKAGE).v$(PREVIOUS_VERSION)
CSV_FILENAME := $(PACKAGE).clusterserviceversion.yaml
CSV_PATH := $(MANIFESTS)/bases/$(CSV_FILENAME)

TMPL := tmpl
YAML := yaml
GEN := gen
NAMESPACE_VAR := {NAMESPACE}
TAG_VAR := {TAG}
IMAGE_VAR := {IMAGE}
VERSION_VAR := {VERSION}
