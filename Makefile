.PHONY: test test-verbose test-coverage clean build

# Run all tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Generate detailed coverage report
test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean up test artifacts
clean:
	rm -f gotask_test coverage.out coverage.html

# Build the binary
build:
	go build -o gotask

# Run a quick test cycle (build + test)
check: build test

# Install the binary
install:
	go install