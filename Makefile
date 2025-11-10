.PHONY: fmt lint test build cover ci tools

# Maximum line length for golines; override with: make fmt LINE_LENGTH=100
LINE_LENGTH ?= 120

fmt:
	@echo "Formatting with gofumpt and golines (max len $(LINE_LENGTH))"
	@command -v gofumpt >/dev/null 2>&1 || { echo "gofumpt not found. Run 'make tools' to install."; exit 1; }
	@command -v golines >/dev/null 2>&1 || { echo "golines not found. Run 'make tools' to install."; exit 1; }
	@gofumpt -l -w .
	@golines --max-len=$(LINE_LENGTH) --base-formatter=gofumpt -w .

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

tools:
	@echo "Installing dev tools: gofumpt and golines"
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/segmentio/golines@latest
