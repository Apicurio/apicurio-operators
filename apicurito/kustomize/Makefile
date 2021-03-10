#
# Make Options
#
MK_OPTIONS := -s

#
# Vars
#

#
# Vars that can be overridden by external env vars
#
DEBUG ?= false
KUBE_USER ?= developer
NAMESPACE ?= apicurito
DEFAULT_IMAGE ?=  quay.io/apicurio/apicurito-operator
IMAGE ?= $(DEFAULT_IMAGE)
DEFAULT_TAG := latest
TAG ?= $(DEFAULT_TAG)

.PHONY: kustomize setup operator app product

kubectl:
ifeq (, $(shell which kubectl))
$(error "No kubectl found in PATH. Please install and re-run")
endif

#
# Setup the installation by installing crds, roles and granting
# privileges for the installing user.
#
setup: kubectl
	$(MAKE) $(MK_OPTIONS) -C bases/crd init
	$(MAKE) $(MK_OPTIONS) -C bases/role init
	#@ Must be invoked by a user with cluster-admin privileges
ifeq ($(DEBUG), false)
	kubectl apply -k setup
else
	kubectl kustomize setup
endif

#
# Install the operator deployment and related resources
#
operator: kubectl
	$(MAKE) $(MK_OPTIONS) -C bases/deployment init
	#@ Can be invoked by a user with namespace privileges (rather than a cluster-admin)
ifeq ($(DEBUG), false)
	kubectl apply -k operator
else
	kubectl kustomize operator
endif

#
# Installs the operator deployment and in addition installs a default CR
#
app: kubectl
	$(MAKE) $(MK_OPTIONS) -C bases/deployment init
	$(MAKE) $(MK_OPTIONS) -C app init
	#@ Can be invoked by a user with namespace privileges (rather than a cluster-admin)
ifeq ($(DEBUG), false)
	kubectl apply -k app
else
	kubectl kustomize app
endif
