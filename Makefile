.PHONY: lint test vendor clean

# Ensure Go modules are used
export GO111MODULE=on

# Default target to run lint and tests
default: lint test

# Lint the code using golangci-lint
lint:
	golangci-lint run

# Run tests
test:
	go test -v -cover ./...

# Run tests using Yaegi (if applicable)
yaegi_test:
	yaegi test -v .

# Vendor dependencies
vendor:
	go mod vendor

# Clean the vendor directory
clean:
	rm -rf ./vendor
