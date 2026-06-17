SHELL = /bin/bash
OS = $(shell uname -s | tr '[:upper:]' '[:lower:]')

GOBIN ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell go env GOPATH)/bin
endif

export PATH := $(GOBIN):$(PATH)

# golangci-lint — pinned version; override with GOLANGCI_LINT_VERSION=vX.Y.Z
GOLANGCI_LINT_VERSION ?= v2.12.2
GOLANGCI_LINT         := $(GOBIN)/golangci-lint

# Security tooling — gosec runs outside golangci-lint (edgelet pattern)
GOVULNCHECK_VERSION ?= v1.1.4
GOSEC_VERSION       ?= v2.22.2
GOSEC_SCOPE         := ./pkg/...

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

# Import path for pkg/apps (must match go.mod module)
APPS_IMPORT_PATH = github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/apps
MODULE_PATH      = github.com/eclipse-iofog/iofog-go-sdk/v3

.PHONY: gen
gen: install-tools ## Generate code
	@sed -i '' -E 's|//(.*// \+k8s:deepcopy-gen=ignore)|\1|g' pkg/apps/types.go
	@sed -i '' -E 's|(.*// \+k8s:deepcopy-gen=ignore)|//\1|g' pkg/apps/types.go
	deepcopy-gen \
		--bounding-dirs $(MODULE_PATH) \
		-i $(APPS_IMPORT_PATH) \
		-O deepcopy_generated \
		-o . \
		-p $(APPS_IMPORT_PATH) \
		--trim-path-prefix $(MODULE_PATH) \
		--go-header-file ./boilerplate.go.txt
	@sed -i '' -E 's|//(.*// \+k8s:deepcopy-gen=ignore)|\1|g' pkg/apps/types.go

$(GOLANGCI_LINT):
	@echo "⬇️  Installing golangci-lint $(GOLANGCI_LINT_VERSION) → $(GOBIN)..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $(GOBIN) $(GOLANGCI_LINT_VERSION)
	@echo "✓ golangci-lint $(GOLANGCI_LINT_VERSION) installed"

.PHONY: install-lint
install-lint: $(GOLANGCI_LINT) ## Install golangci-lint (pinned to GOLANGCI_LINT_VERSION)
	@$(GOLANGCI_LINT) version

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Run linters (auto-installs golangci-lint if needed)
	@echo "Running golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@$(GOLANGCI_LINT) run --config .golangci.yaml ./pkg/...

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT) ## Run linters and auto-fix issues where possible
	@echo "Running golangci-lint $(GOLANGCI_LINT_VERSION) with --fix..."
	@$(GOLANGCI_LINT) run --config .golangci.yaml ./pkg/... --fix

.PHONY: fmt
fmt: ## Format the source
	@gofmt -s -w .

.PHONY: security-code
security-code: ## Static Go security analysis (gosec; not in golangci-lint)
	@echo "🔍 Running Go static security analysis..."
	@if ! command -v gosec >/dev/null 2>&1; then \
		echo "⬇️  Installing gosec $(GOSEC_VERSION)..."; \
		go install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION); \
	fi
	@gosec -exclude-dir=reports -exclude-generated $(GOSEC_SCOPE)

.PHONY: vulncheck
vulncheck: ## Dependency vulnerability scan (govulncheck + go mod verify)
	@echo "🔐 Running govulncheck..."
	@if ! command -v govulncheck >/dev/null 2>&1; then \
		echo "⬇️  Installing govulncheck $(GOVULNCHECK_VERSION)..."; \
		go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION); \
	fi
	@govulncheck ./pkg/...
	@echo "🔍 Verifying module integrity..."
	@go mod verify

.PHONY: quality
quality: lint security-code vulncheck ## Full Plan 4 quality gate

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


# Pin code-generator to a version compatible with Go 1.21+ (avoid @latest which requires Go 1.25+)
DEEPCOPY_GEN_VERSION ?= v0.29.0

.PHONY: install-tools
install-tools:
	go install -v k8s.io/code-generator/cmd/deepcopy-gen@$(DEEPCOPY_GEN_VERSION)
