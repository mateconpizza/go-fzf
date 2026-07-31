.PHONY: all clean test testcover testsum

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
