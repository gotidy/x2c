SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

mod:
	go mod tidy
.PHONY: mod

test: mod
	go test ./...
.PHONY: test

## Format the Code.
fmt:
	// gofmt -s -l -w $(SRCS)
	gofumpt -s -l -w $(SRCS)
.PHONY: fmt

## Lint the Code.
lint: mod
	golangci-lint run -v --out-format=tab --timeout 10m0s
.PHONY: lint

## Install tools
tools:
	go install github.com/mvdan/gofumpt@latest
	go install github.com/goreleaser/goreleaser@latest
	if [[ "$$OSTYPE" == "darwin"* ]]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.43.0; fi	
	if [[ "$$OSTYPE" == "linux-gnu"* ]]; then wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.43.0; fi
.PHONY: tools

pre-commit:
	pip3 install pre-commit
	pre-commit install
.PHONY: pre-commit

release: 
	goreleaser release --rm-dist
.PHONY: release

snapshot: .goreleaser.yml tools/goreleaser
	goreleaser --snapshot --rm-dist
.PHONY: snapshot
	
