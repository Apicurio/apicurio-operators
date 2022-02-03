
CONFIG := config
#
# Include the main common Makefile containing
# basic common recipes and vars
#
include $(CONFIG)/common/Makefile

# kernel-style V=1 build verbosity
ifeq ("$(origin V)", "command line")
       BUILD_VERBOSE = $(V)
endif

.PHONY: all
all: build

.PHONY: format
format:
	./scripts/go-fmt.sh

.PHONY: vet
vet:
	./scripts/go-vet.sh

.PHONY: test
test:
	./scripts/go-test.sh

.PHONY: lint
lint:
	# Temporarily disabled
	# ./scripts/go-lint.sh
	# ./scripts/yaml-lint.sh

.PHONY: build
build:
	IMAGE=$(IMAGE) TAG=$(TAG) VERSION=$(VERSION) ./scripts/go-build.sh


.PHONY: clean
clean:
	rm -rf build/_output

.PHONY: setup
setup:
	$(MAKE) $(MK_OPTIONS) -C $(CONFIG) $@

.PHONY: operator
operator:
	$(MAKE) $(MK_OPTIONS) -C $(CONFIG) $@

.PHONY: app
app:
	$(MAKE) $(MK_OPTIONS) -C $(CONFIG) $@
