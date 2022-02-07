#
# Allows for resources to be loaded from outside the root location of
# the kustomize config file. Ensures that resource don't need to be
# copied around the file system.
#
# See https://kubectl.docs.kubernetes.io/faq/kustomize
#
KOPTIONS := --load-restrictor LoadRestrictionsNone

#
# Include the main common Makefile containing
# basic common recipes and vars
#
include common/Makefile

#
# Vars that can be overridden by external env vars
#
DEBUG ?= false
KUBE_USER ?= developer
LEGACY ?= false

NAMESPACE := $(shell ./script/namespace.sh)
PACKAGE := apicurito
DOMAIN := apicur.io
BUNDLE_DIR := bundle
CATALOG := catalog
PKG := ../pkg
CRD := crd/bases
MANAGER := manager
MANIFESTS := manifests
RBAC := rbac

CSV_SUPPORT := Apicurito
CSV_DISPLAY_NAME := Apicurito
CSV_NAME := $(PACKAGE).v$(VERSION)
CSV_REPLACES := $(PACKAGE).v$(PREVIOUS_VERSION)
CSV_FILENAME := $(PACKAGE).clusterserviceversion.yaml
CSV_PATH := $(MANIFESTS)/bases/$(CSV_FILENAME)

.PHONY: kustomize operator-sdk setup operator app

kubectl:
ifeq (, $(shell which kubectl))
$(error "No kubectl found in PATH. Please install and re-run")
endif

#
# Setup the installation by installing crds, roles and granting
# privileges for the installing user.
#
setup: kubectl
	$(MAKE) $(MK_OPTIONS) -C $(RBAC) init
	#@ Must be invoked by a user with cluster-admin privileges
ifeq ($(LEGACY), true)
	@cd setup && \
		$(KUSTOMIZE) edit remove resource ../crd/bases && \
		$(KUSTOMIZE) edit add resource ../crd/deprecated
else
	@cd setup && \
		$(KUSTOMIZE) edit remove resource ../crd/deprecated && \
		$(KUSTOMIZE) edit add resource ../crd/bases
endif
ifeq ($(DEBUG), false)
	@$(KUSTOMIZE) build $(KOPTIONS) setup | kubectl apply -f -
else
	@$(KUSTOMIZE) build $(KOPTIONS) setup
endif

#
# Install the operator deployment and related resources
#
operator: kubectl
	$(MAKE) $(MK_OPTIONS) -C $(MANAGER) init
	#@ Can be invoked by a user with namespace privileges (rather than a cluster-admin)
ifeq ($(DEBUG), false)
	@$(KUSTOMIZE) build $(KOPTIONS) operator | kubectl apply -f -
else
	@$(KUSTOMIZE) build $(KOPTIONS) operator
endif

#
# Installs the operator deployment and in addition installs a default CR
#
app: kubectl
	$(MAKE) $(MK_OPTIONS) -C $(MANAGER) init
	#@ Can be invoked by a user with namespace privileges (rather than a cluster-admin)
ifeq ($(DEBUG), false)
	@$(KUSTOMIZE) build $(KOPTIONS) app | kubectl apply -f -
else
	@$(KUSTOMIZE) build $(KOPTIONS) app
endif


#
##################################
# B U N D L E  G E N E R A T I O N
##################################
#
# Default bundle image tag
BUNDLE_IMG ?= $(IMAGE)-bundle:$(VERSION)

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

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION) ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

kustomize:
ifeq (, $(shell which kustomize))
	@{ \
	set -e ;\
	KUSTOMIZE_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$KUSTOMIZE_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/kustomize/kustomize/v4@$(KUSTOMIZE_VERSION) ;\
	rm -rf $$KUSTOMIZE_GEN_TMP_DIR ;\
	}
KUSTOMIZE=$(GOBIN)/kustomize
else
KUSTOMIZE=$(shell which kustomize)
endif

opm:
ifeq (, $(shell which opm))
	@{ \
	set -e ;\
	OPM_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$OPM_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get \
		-ldflags '-w -extldflags "-static"' -tags "json1" \
		-ldflags "-X '$(OPM_PKG)/cmd/opm/version.opmVersion=$(OPM_VERSION)' -X '$(OPM_PKG)/cmd/opm/version.buildDate=$(BUILD_TIME)'" \
		$(OPM_PKG)/cmd/opm@$(OPM_VERSION) ;\
	rm -rf $$OPM_GEN_TMP_DIR ;\
	}
OPM=$(GOBIN)/opm
else
OPM=$(shell which opm)
endif

operator-sdk:
ifeq (, $(shell which operator-sdk))
	@{ \
	set -e ;\
	curl \
		-L https://github.com/operator-framework/operator-sdk/releases/download/$(OPERATOR_SDK_VERSION)/operator-sdk_linux_amd64 \
		-o operator-sdk ;\
	chmod +x operator-sdk ;\
	mv operator-sdk $(GOBIN)/ ;\
	}
OPERATOR_SDK=$(GOBIN)/operator-sdk
else
OPERATOR_SDK=$(shell which operator-sdk)
endif

.PHONY: bundle manifests generate generate-deepcopy generate-crds

generate-deepcopy: controller-gen
	@$(CONTROLLER_GEN) paths="$(PKG)/..." object

generate-crds: controller-gen
	@$(CONTROLLER_GEN) crd crd:crdVersions=v1 \
		paths=$(PKG)/apis/... \
		output:artifacts:config=$(CRD) \
		output:crd:dir=$(CRD)
	@sed -i '/creationTimestamp: null/d' $(CRD)/$(DOMAIN)_*.yaml

# Generate code
generate: generate-deepcopy generate-crds

CSV=$(MANIFESTS)/bases/apicurito.clusterserviceversion.yaml
TIMESTAMP=$(shell date -u '+%F %T %Z')

# Tailor the manifest according to default values for this project
# Note. to make the bundle this name must match that specified in PROJECT
#
pre-bundle:
# bundle name must match that which appears in PROJECT file
	@sed -i 's/projectName: .*/projectName: $(PACKAGE)/' PROJECT
# finds the single CSV file and renames it
	@find $(MANIFESTS)/bases -type f \( -name "*.clusterserviceversion.yaml" ! -name $(CSV_FILENAME) \) -execdir mv '{}' $(CSV_FILENAME) ';'
	@sed -i 's~^    containerImage: .*~    containerImage: $(IMAGE):$(TAG)~' $(CSV_PATH)
	@sed -i 's/^    support: .*/    support: $(CSV_SUPPORT)/' $(CSV_PATH)
	@sed -i 's/^  name: .*.\(v.*\)/  name: $(CSV_NAME)/' $(CSV_PATH)
	@sed -i 's/^      displayName: .*/      displayName: $(CSV_DISPLAY_NAME)/' $(CSV_PATH)
	@sed -i 's/^  replaces: .*/  replaces: $(CSV_REPLACES)/' $(CSV_PATH)
	@sed -i 's/^  version: .*/  version: $(VERSION)/' $(CSV_PATH)

# Generate bundle manifests and metadata, then validate generated files.
bundle: generate-crds kustomize pre-bundle
	@$(MAKE) $(MK_OPTIONS) -C $(RBAC) init
	@$(MAKE) $(MK_OPTIONS) -C $(MANAGER) init
# (Re)-generate the manifests directory
	@$(OPERATOR_SDK) generate kustomize manifests \
		--apis-dir $(PKG)/apis \
		--input-dir $(MANIFESTS) \
		--output-dir $(MANIFESTS)
# Sets the createdAt timestamp
	@sed -i 's/createdAt: .*/createdAt: $(TIMESTAMP)/' $(CSV)
# Sets the operator image to the preferred image:tag
	@cd $(MANIFESTS) && $(KUSTOMIZE) edit set image $(DEFAULT_IMAGE)=$(IMAGE):$(TAG)
# Generates the bundle complete with manifests
	@$(KUSTOMIZE) build $(KOPTIONS) $(MANIFESTS) | \
		$(OPERATOR_SDK) generate bundle \
		-q --overwrite --version $(VERSION) \
		--kustomize-dir $(MANIFESTS) $(BUNDLE_METADATA_OPTS)
# # Moves the docker file into the bundle directory
	@mv bundle.Dockerfile $(BUNDLE_DIR)/Dockerfile && \
		sed -i 's/bundle\///g' $(BUNDLE_DIR)/Dockerfile
# Add Product LABELS to Dockerfile
	@cat Dockerfile.labels >> $(BUNDLE_DIR)/Dockerfile
	@$(OPERATOR_SDK) bundle validate $(BUNDLE_DIR)

# Build the bundle image.
.PHONY: bundle-build
bundle-build: bundle
	cd bundle && docker build -f Dockerfile -t $(BUNDLE_IMG) .

.PHONY: catalog
catalog: opm
	SRC_CATALOG=$(SRC_CATALOG) CATALOG_DIR=$(CATALOG) \
	OPM=$(OPM) BUNDLE_IMAGE=$(BUNDLE_IMG) CSV_NAME=$(CSV_NAME) \
	CSV_REPLACES=$(CSV_REPLACES) CHANNEL=$(DEFAULT_CHANNEL) ./script/build-catalog.sh

clean:
	find . -name "*.gen.*" -delete
	rm -rf bundle
