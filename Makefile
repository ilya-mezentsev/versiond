ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

APP_DIR := $(ROOT_DIR)/app
TESTS_DIR := $(ROOT_DIR)/tests

LINTER_DIR := $(ROOT_DIR)/linter
LINTER_BIN := $(LINTER_DIR)/bin/golangci-lint

ENTRYPOINT := $(APP_DIR)/cmd/main.go
CONFIG := $(APP_DIR)/configs/main.json

TEST_SERVER_ENTRYPOINT := $(TESTS_DIR)/main.go

setup: tidy mocks test-scripts install-linter

tidy:
	go mod tidy

run:
	CONFIG_PATH=$(CONFIG) go run $(ENTRYPOINT)

mocks:
	cd $(APP_DIR) && mockery

install-linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LINTER_DIR)/bin v1.57.2

# about GOEXPERIMENT see https://github.com/golang/go/issues/65570
test:
	cd $(APP_DIR) && GOEXPERIMENT=nocoverageredesign go test -race -cover ./internal/... | { grep -v "no test files"; true; }

lint:
	cd $(APP_DIR) && $(LINTER_BIN) run ./internal/...

build:
	go build -o $(ROOT_DIR)/built/main $(ENTRYPOINT)

test-server:
	go run $(TEST_SERVER_ENTRYPOINT)

test-scripts:
	cp $(TESTS_DIR)/*.sh /tmp/
