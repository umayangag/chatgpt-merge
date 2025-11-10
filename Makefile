.PHONY: fmt lint test build cover ci

fmt:
	@gofmt -s -w .

lint:
	@go vet ./...

test:
	@go test ./...

build:
	@go build -o chatgpt-merge ./cmd

cover:
	@go test -coverprofile=profile.cov ./...
	@go tool cover -func=profile.cov

ci: fmt lint test
