ASSETS := ../../../deploy
TMPL := tmpl
YAML := yaml
NAMESPACE_VAR := {NAMESPACE}
USER_VAR := {KUBE_USER}

CLUSTER_ROLE := ./cluster_role.gen
CLUSTER_ROLE_BINDING := ./cluster_role_binding.gen
ROLE := ./role.gen
ROLE_BINDING := ./role_binding.gen
ROLE_USER := ./role_app_user.yaml

NAMESPACE ?= apicurito

.PHONY: sync init

#
# start-sync
#
# Copy the template from the src directory
# Convert the template to a formatted yaml file:
#
sync:
	cp $(ASSETS)/cluster_role.yaml $(CLUSTER_ROLE).$(TMPL)
	cp $(ASSETS)/cluster_role_binding.yaml $(CLUSTER_ROLE_BINDING).$(TMPL)
	sed -i 's/namespace: apicurito/namespace: $(NAMESPACE_VAR)/' $(CLUSTER_ROLE_BINDING).$(TMPL)
	sed -i '/metadata:/,/roleRef:/s/apicurito/apicurito-$(NAMESPACE_VAR)/' $(CLUSTER_ROLE_BINDING).$(TMPL)
	cp $(ASSETS)/role.yaml $(ROLE).$(TMPL)
	cp $(ASSETS)/role_binding.yaml $(ROLE_BINDING).$(TMPL)
# end-sync

init: sync
	for resource in $(CLUSTER_ROLE) $(CLUSTER_ROLE_BINDING) $(ROLE) $(ROLE_BINDING); do \
		cp $${resource}.$(TMPL) $${resource}.$(YAML); \
		sed -i 's/$(NAMESPACE_VAR)/$(NAMESPACE)/' $${resource}.$(YAML); \
	done
