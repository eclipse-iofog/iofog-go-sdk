SHELL = /bin/bash
OS = $(shell uname -s | tr '[:upper:]' '[:lower:]')

# Build variables
VERSION ?= $(shell git tag | tail -1 | sed "s|v||g")-dev
COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)
REPORTS_DIR ?= reports
TEST_RESULTS ?= TEST-iofog-go-sdk.txt
TEST_REPORT ?= TEST-iofog-go-sdk.xml

.PHONY: init
init: ## Init git repository
	@cp gitHooks/* .git/hooks/

.PHONY: all
all: test ## Generate code and run tests

.PHONY: clean
clean: ## Clean the working area and the project
	rm -rf $(REPORTS_DIR)

.PHONY: gen
gen: install-tools ## Generate code
	@sed -i'' -E "s|//(.*// \+k8s:deepcopy-gen=ignore)|\1|g" pkg/apps/types.go
	@sed -i'' -E "s|(.*// \+k8s:deepcopy-gen=ignore)|//\1|g" pkg/apps/types.go
	deepcopy-gen -i ./pkg/apps -o . --go-header-file ./boilerplate.go.txt
	@sed -i'' -E "s|//(.*// \+k8s:deepcopy-gen=ignore)|\1|g" pkg/apps/types.go

.PHONY: lint
lint: golangci-lint fmt ## Lint the source
	@$(GOLANGCI_LINT) run --timeout 5m0s

golangci-lint: ## Install golangci
ifeq (, $(shell which golangci-lint))
	@{ \
	set -e ;\
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1 ;\
	}
GOLANGCI_LINT=$(GOBIN)/golangci-lint
else
GOLANGCI_LINT=$(shell which golangci-lint)
endif

.PHONY: fmt
fmt: ## Format the source
	@gofmt -s -w .

.PHONY: test
test: gen fmt ## Run unit tests
	mkdir -p $(REPORTS_DIR)
	rm -f $(REPORTS_DIR)/*
	set -o pipefail; go list ./pkg/... | xargs -n1 go test -ldflags "$(LDFLAGS)" -v -parallel 1 2>&1 | tee $(REPORTS_DIR)/$(TEST_RESULTS)

.PHONY: list
list: ## List all make targets
	@$(MAKE) -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: help
.DEFAULT_GOAL := help
help: ## Get help output
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Variable outputting/exporting rules
var-%: ; @echo $($*)
varexport-%: ; @echo $*=$($*)


.PHONE: install-tools
install-tools:
	env | grep GO
	env | grep PATH
	go install -v k8s.io/code-generator/cmd/deepcopy-gen@v0.26
	which deepcopy-gen
