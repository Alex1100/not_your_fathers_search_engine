#!make
include .env
export $(shell sed 's/=.//' .env)

.PHONY: test ci-check

lint: lint-check-deps
	@echo "[golangci-lint] linting sources"
	@golangci-lint run \
		-E misspell \
		-E golint \
		-E gofmt \
		-E unconvert \
		--exclude-use-default=false \
		./...

test: 
	@echo "[go test] running tests and collecting coverage metrics"
	@go test -v -tags all_tests -race -coverprofile=coverage.txt -covermode=atomic ./...

lint-check-deps:
	@if [ -z `which golangci-lint` ]; then \
			echo "[go get] installing golangci-lint";\
			curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.27.0;\
	fi

ci-check: lint test

commit: lint
	@echo "commiting code"
	@git send "$(m)"