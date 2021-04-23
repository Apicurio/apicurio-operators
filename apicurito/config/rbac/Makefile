TMPL := tmpl
YAML := yaml
GEN := gen
NAMESPACE_VAR := {NAMESPACE}

CLUSTER_ROLE_BINDING := ./cluster_role_binding

NAMESPACE := $(shell ../script/namespace.sh)

.PHONY: init

init:
	for resource in $(CLUSTER_ROLE_BINDING); do \
		cp $${resource}.$(TMPL) $${resource}.$(GEN).$(YAML); \
		sed -i 's/$(NAMESPACE_VAR)/$(NAMESPACE)/' $${resource}.$(GEN).$(YAML); \
	done
