
#
# Include the main common Makefile containing
# basic common recipes and vars
#
include config/common/Makefile

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
	IMAGE=$(IMAGE) TAG=$(TAG) ./scripts/go-build.sh


.PHONY: clean
clean:
	rm -rf build/_output
