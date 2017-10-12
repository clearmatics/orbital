# Copyright (C) 2017 Clearmatics - All Rights Reserved

# This Makefile is meant to be used mainly by CI
.PHONY: all test clean

all:
	build/env.sh go get
	build/env.sh go build

test: 	all
	build/env.sh go get github.com/stretchr/testify
	build/env.sh go test ./...

coverage: 	all
	build/env.sh go get github.com/stretchr/testify
	build/env.sh build/coverage.sh
	build/env.sh go tool cover -html=coverage.out -o=coverage.html

clean:
	rm -fr build/_workspace/Godeps/

format:
	build/env.sh gofmt -s -w .

