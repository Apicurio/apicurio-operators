TMPL := tmpl
YAML := yaml
GEN := gen
NAMESPACE_VAR := {NAMESPACE}
TAG_VAR := {TAG}
IMAGE_VAR := {IMAGE}

DEPLOYMENT := ./deployment

#
# Include the main common Makefile containing
# basic common recipes and vars
#
include ../common/Makefile

NAMESPACE := $(shell ../script/namespace.sh)

.PHONY: init

init:
	cp $(DEPLOYMENT).$(TMPL) $(DEPLOYMENT).$(GEN).$(YAML)
	sed -i 's/$(NAMESPACE_VAR)/$(NAMESPACE)/' $(DEPLOYMENT).$(GEN).$(YAML)
	sed -i 's/$(TAG_VAR)/$(TAG)/' $(DEPLOYMENT).$(GEN).$(YAML)
	sed -i 's~$(IMAGE_VAR)~$(IMAGE)~' $(DEPLOYMENT).$(GEN).$(YAML)
