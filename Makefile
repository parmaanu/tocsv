REPO_NAME := "tocsv"
PKG_LIST := $(shell go list ${REPO_NAME}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep build clean test coverage lint

all: build

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	@CC=clang CXX=clang++ go test -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage report (also in HTML)
	./scripts/coverage.sh .;

dep:
	@go get -v -d ./...

build: dep
	@go build -o tocsv.out .

# clean: ## Remove previous build


help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
