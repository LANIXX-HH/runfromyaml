# build the Go binary
build:
	go build
	cp runfromyaml ${HOME}/bin/

# run all tests
test:
	go test -v ./...

# run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# run tests with race detection
test-race:
	go test -v -race ./...

# run benchmarks
benchmark:
	go test -v -bench=. ./...

# run tests and generate coverage report
test-full: test-coverage test-race
	@echo "Full test suite completed. Coverage report available at coverage.html"

# run specific package tests
test-config:
	go test -v ./pkg/config/...

test-cli:
	go test -v ./pkg/cli/...

test-functions:
	go test -v ./pkg/functions/...

test-errors:
	go test -v ./pkg/errors/...

# install dependencies
deps:
	go get -v ./...
	go mod download

# update to latest
update:
	go get -u ./...

# clean the project
clean:
	rm -rf runfromyaml
	rm -rf ${HOME}/bin/runfromyaml
	rm -rf coverage.out coverage.html

# lint the code (requires golangci-lint)
lint:
	golangci-lint run

# format the code
fmt:
	go fmt ./...

# vet the code
vet:
	go vet ./...

# run all quality checks
quality: fmt vet lint test

# setup pre-commit hooks
setup-precommit:
	./scripts/setup-precommit.sh

# run pre-commit on all files
precommit-all:
	pre-commit run --all-files

.PHONY: build test test-coverage test-race benchmark test-full test-config test-cli test-functions test-errors deps update clean lint fmt vet quality setup-precommit precommit-all
