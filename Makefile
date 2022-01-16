.PHONY: mod test fmt lint tools pre-commit

SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

mod:
	go mod tidy

test: mod
	go test ./...

## Format the Code.
fmt:
	// gofmt -s -l -w $(SRCS)
	gofumpt -s -l -w $(SRCS)

## Lint the Code.
lint: mod
	golangci-lint run -v --out-format=tab --timeout 10m0s

## Install tools
tools:
	go install github.com/mvdan/gofumpt
	if [[ "$$OSTYPE" == "darwin"* ]]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.43.0; fi	
	if [[ "$$OSTYPE" == "linux-gnu"* ]]; then wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.43.0; fi

pre-commit:
	pip3 install pre-commit
	pre-commit install
