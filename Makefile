GO       := go
GOTEST   := $(GO) test
COVERAGE := coverage.out
HTML     := coverage.html
FN	 ?= .

all: test

test:
	@echo ">> Running tests"
	@$(GOTEST) -run="^Test" -v ./...

testcover:
	@echo ">> Running tests with coverage"
	@$(GOTEST) ./... -coverprofile=$(COVERAGE)
	@$(GO) tool cover -html=$(COVERAGE) -o $(HTML)
	@xdg-open $(HTML)

testsum:
	@echo ">> Running tests with gotestsum"
	@gotestsum --format=github-actions --hide-summary=skipped

example:
	@echo ">> Running test examples"
	@go test -run="^Example"

# Run tests for a specific function
testfn:
	@echo '>> Testing function $(FN)'
	@go test -run $(FN) ./...

clean:
	@echo ">> Cleaning"
	rm -f $(COVERAGE) $(HTML)

ci:
	@if go list -m -f '{{if .Replace}}{{if not .Replace.Version}}{{.Path}} => {{.Replace.Path}}{{end}}{{end}}' all | grep -q .; then \
		echo "error: local replace directive found in go.mod"; \
		go list -m -f '{{if .Replace}}{{if not .Replace.Version}}{{.Path}} => {{.Replace.Path}}{{end}}{{end}}' all; \
		exit 1; \
	fi
	go mod tidy
	git diff --exit-code
	go vet ./...
	go build ./...
	@$(GOTEST) -run="^Test" -race ./...

lint:
	@echo '>> Linting code'
	@go vet ./...
	golangci-lint run ./...

.PHONY: all test clean lint testfn testcover testsum
