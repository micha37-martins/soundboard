.DEFAULT_GOAL := help

GOLANGCI_LINT_VERSION ?= v1.29.0 ## ?= means: if not set or empty set to v...
TEST_FLAGS ?= -race ## add race condition test flag
PKGS ?= $(shell go list ./... | grep -v /vendor/) # list packages folder path

## help lists all available commands from this Makefile
.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "[32m%-23s[0m %s\n", $$1, $$2}'

.PHONY: binary
binary: deps build ## install binary dependencies and build

## https://dev.to/defman/introducing-go-mod-1cdo
.PHONY: deps
deps: ## install go deps
	go mod download

## ldflags add dynamic informations into binary (https://golang.org/cmd/link/)
.PHONY: build
build: ## build soundboard
	go build -ldflags="-s -w" -o soundboard cmd/main.go

.PHONY: checkup
checkup: lint test vet coverage ## checks for errors and bad code

.PHONY: test
test: ## run tests
	go test $(TEST_FLAGS) $(PKGS)

.PHONY: vet
vet: ## run go vet to check code for suspicious constructs
	go vet $(PKGS)

.PHONY: coverage
coverage: ## generate code coverage overview
	go test $(TEST_FLAGS) -covermode=atomic -coverprofile=coverage.txt $(PKGS)
	go tool cover -func=coverage.txt

.PHONY: lint
lint: ## run golangci-lint
	command -v golangci-lint > /dev/null 2>&1 || \
	  curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)
	golangci-lint run

.PHONY: clean
clean: ## clean dependencies and artifacts
	rm -rf vendor/
	rm -f soundboard

.PHONY: install
install: build ## install soundboard into $GOPATH/bin
	mv soundboard $(GOPATH)/bin/soundboard
