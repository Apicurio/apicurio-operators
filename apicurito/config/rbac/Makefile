CONFIG_DIR := ../
# User customisable variables
include $(CONFIG_DIR)/vars/Makefile

CLUSTER_ROLE_BINDING := ./cluster_role_binding

.PHONY: init

init:
	for resource in $(CLUSTER_ROLE_BINDING); do \
		cp $${resource}.$(TMPL) $${resource}.$(GEN).$(YAML); \
		sed -i 's/$(NAMESPACE_VAR)/$(NAMESPACE)/' $${resource}.$(GEN).$(YAML); \
	done
