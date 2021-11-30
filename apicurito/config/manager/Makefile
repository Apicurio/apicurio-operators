CONFIG_DIR := ../
# User customisable variables
include $(CONFIG_DIR)/vars/Makefile

DEPLOYMENT := ./deployment

.PHONY: init

init:
	cp $(DEPLOYMENT).$(TMPL) $(DEPLOYMENT).$(GEN).$(YAML)
	sed -i 's/$(NAMESPACE_VAR)/$(NAMESPACE)/' $(DEPLOYMENT).$(GEN).$(YAML)
	sed -i 's/$(TAG_VAR)/$(TAG)/' $(DEPLOYMENT).$(GEN).$(YAML)
	sed -i 's~$(IMAGE_VAR)~$(IMAGE)~' $(DEPLOYMENT).$(GEN).$(YAML)
	sed -i 's~$(VERSION_VAR)~$(VERSION)~' $(DEPLOYMENT).$(GEN).$(YAML)
