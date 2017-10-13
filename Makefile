GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
SHELL=/bin/bash

test:
	go test ./...

format:
	gofmt -s -w .

coverage:
	go tool cover -html=coverage.out -o=coverage.html

build:
	go build .

check:
	@if [ -n "$(shell gofmt -l ${GOFILES_NOVENDOR})" ]; then \
		echo 1>&2 'The following files need to be formatted:'; \
		gofmt -l .; \
		exit 1; \
		fi

vet:
	@go vet ${GOFILES_NOVENDOR}

lint:
	golint ${GOFILES_NOVENDOR}



