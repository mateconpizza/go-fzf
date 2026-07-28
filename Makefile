.PHONY: all clean test testcover testsum

GO       := go
GOTEST   := $(GO) test
COVERAGE := coverage.out
HTML     := coverage.html

all: test

test:
	@echo ">> Running tests"
	@$(GOTEST) -v ./...

testcover:
	@echo ">> Running tests with coverage"
	@$(GOTEST) ./... -coverprofile=$(COVERAGE)
	@$(GO) tool cover -html=$(COVERAGE) -o $(HTML)
	@xdg-open $(HTML)

testsum:
	@echo ">> Running tests with gotestsum"
	@gotestsum --format=github-actions --hide-summary=skipped

clean:
	@echo ">> Cleaning"
	rm -f $(COVERAGE) $(HTML)
