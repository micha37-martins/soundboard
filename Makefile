.DEFAULT_GOAL := help

GOLANGCI_LINT_VERSION ?= v1.29.0 # if not set or empty set to v...
TEST_FLAGS ?= -race # add race condition test flag
PKGS ?= $(shell go list ./... | grep -v /vendor/) # list packages folder path

# help lists all available commands from this Makefile
.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "[32m%-23s[0m %s\n", $$1, $$2}'

.PHONY: all
all: app binary ## install dependencies and build everything

.PHONY: binary
binary: deps pack-app build ## install binary dependencies, pack app and build

.PHONY: deps
deps: ## install go deps
  # https://dev.to/defman/introducing-go-mod-1cdo
	go mod download

.PHONY: build
build: ## build soundboard
	# ldflags add dynamic informations into binary (https://golang.org/cmd/link/)
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

.PHONY: images
images: image-amd64 image-armv7 ## build docker images

.PHONY: image-amd64
image-amd64: ## build amd64 image
	docker build --build-arg GOARCH=amd64 -t micha37martins/soundboard:amd64 .

.PHONY: image-armv7
image-armv7: ## build armv7 image
	docker build --build-arg GOARCH=arm --build-arg GOARM=7 -t micha37martins/soundboard:armv7 .
