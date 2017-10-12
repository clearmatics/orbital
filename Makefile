.PHONY: all test clean

test: 	all
	go test ./...

format: all
	gofmt -s -w .

coverage: 	all
	go tool cover -html=coverage.out -o=coverage.html


