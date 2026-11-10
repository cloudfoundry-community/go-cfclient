GOLANG_CI_LINT_VERSION ?= v2.15.0
GOLANG_CI_LINT := ./bin/golangci-lint

.PHONY: all
all: test lint

.PHONY: clean
clean: ## Clean testcache and delete build output
	go clean -testcache

.PHONY: test
test: ## Run the unit tests
	go test ./... -v -race

.PHONY: generate
generate: ## Generate fakes
	go generate ./...

.PHONY: lint-prepare
lint-prepare:
	@if ! $(GOLANG_CI_LINT) version --short 2>/dev/null | grep -Fxq '$(patsubst v%,%,$(GOLANG_CI_LINT_VERSION))'; then \
		echo "Installing golangci-lint $(GOLANG_CI_LINT_VERSION)"; \
		installer=$$(mktemp) || exit 1; \
		trap 'rm -f "$$installer"' EXIT HUP INT TERM; \
		curl -sSfL https://golangci-lint.run/install.sh -o "$$installer" && \
		sh "$$installer" -b ./bin $(GOLANG_CI_LINT_VERSION) || exit 1; \
	fi

.PHONY: lint
lint: lint-prepare ## Run the golangci linter
	$(GOLANG_CI_LINT) run

.PHONY: update
update: ## Update all dependencies
	go get -u ./... && go mod tidy

.PHONY: tidy
tidy: ## Remove unused dependencies
	go mod tidy

.PHONY: list
list: ## Print the current module's dependencies.
	go list -m all

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print help for each make target
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
