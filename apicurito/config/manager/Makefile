CONFIG_DIR := ../
# User customisable variables
include $(CONFIG_DIR)/vars/Makefile

DEPLOYMENT := ./deployment
# cut out the last version digit
SHORT_VERSION=$(shell echo $(VERSION)|cut -d'.' -f 1,2)

.PHONY: init

init:
	cp $(DEPLOYMENT).$(TMPL) $(DEPLOYMENT).$(GEN).$(YAML)
	sed -i 's/$(NAMESPACE_VAR)/$(NAMESPACE)/' $(DEPLOYMENT).$(GEN).$(YAML)
	sed -i 's/$(TAG_VAR)/$(TAG)/' $(DEPLOYMENT).$(GEN).$(YAML)
	sed -i 's~$(IMAGE_VAR)~$(IMAGE)~' $(DEPLOYMENT).$(GEN).$(YAML)
	sed -i '/rht/s/$(VERSION_VAR)/$(SHORT_VERSION)/' $(DEPLOYMENT).$(GEN).$(YAML)
	sed -i 's~$(VERSION_VAR)~$(VERSION)~' $(DEPLOYMENT).$(GEN).$(YAML)
