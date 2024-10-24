GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")
PACKAGES ?= $(shell $(GO) list ./...)
PROJECT_PATH=$(shell cd "$(dirname "$0")"; pwd)
DIR_CONF ?= /etc/pmon3/config
TEST_REGEX := $(or $(TEST_REGEX),"Test")
TEST_FILE_CONFIG ?= $(PROJECT_PATH)/test/e2e/config/test-config.core.yml
TEST_DIR_LOGS=$(shell cat $(TEST_FILE_CONFIG) | grep "directory:" | sed -n "1 p" | cut -d' ' -f4)
TEST_ARTIFACT_PATH=$(shell dirname "$(TEST_DIR_LOGS)")
DEFAULT_TEST_PACKAGES := "./..."
TEST_PACKAGES := $(or $(TEST_PACKAGES),$(DEFAULT_TEST_PACKAGES))
COVERAGE_OMISSION := '!/^(pmon3\/utils|pmon3\/test|pmon3\/cmd|pmon3\/cli\/cobra|pmon3\/pmond\/protos)/'

all: help

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## remove files created during build pipeline
	$(call print-target)
	rm -f coverage.*
	rm -f '"$(shell go env GOCACHE)/../golangci-lint"'
	go clean -i -cache -testcache -fuzzcache -x
	rm -rf "$(TEST_ARTIFACT_PATH)"

.PHONY: fmt
fmt: ## format files
	$(call print-target)
	$(GOFMT) -w $(GOFILES)

.PHONY: lint
lint: ## lint files
	$(call print-target)
	golangci-lint run --fix

.PHONY: misspell
misspell: ## check for misspellings
	$(call print-target)
	misspell -error $(GOFILES)

.PHONY: betteralign
betteralign: ## check for better aligned structs
	$(call print-target)
	betteralign ./...

.PHONY: tools
tools: ## go install tools
	$(call print-target)
	cd tools && go install $(shell cd tools && $(GO) list -e -f '{{ join .Imports " " }}' -tags=tools)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

.PHONY: mod
mod: ## go mod tidy
	$(call print-target)
	go mod tidy
	cd tools && go mod tidy

.PHONY: base_build
base_build: mod fmt tools misspell betteralign
	cd tools && $(GO) mod tidy
	$(ENV_VARS) $(GO) build $(BUILD_FLAGS) -o bin/pmon3 cmd/pmon3/pmon3.go
	$(ENV_VARS) $(GO) build $(BUILD_FLAGS) -o bin/pmond cmd/pmond/pmond.go

.PHONY: build
build: ENV_VARS=CGO_ENABLED=0
build: base_build ## run tidy and build with CGO disabled
	$(call print-target)

.PHONY: build_cgo
build_cgo: ENV_VARS=CGO_ENABLED=1
build_cgo: base_build  ## run tidy and build with CGO enabled
	$(call print-target)

.PHONY: test
test: build run_test ## build and run tests with CGO disabled
	$(call print-target)

.PHONY: test_cgo
test_cgo: BUILD_FLAGS=$(shell echo '-tags posix_mq,cgo_sqlite')
test_cgo: build_cgo run_test ## build and run tests with CGO enabled
	$(call print-target)

.PHONY: test_net
test_net: BUILD_FLAGS=$(shell echo '-tags net')
test_net: build run_test ## build and run tests with net build tag
	$(call print-target)

.PHONY: make_test_app
make_test_app: ## build the test app
	$(call print-target)
	mkdir -p "$(TEST_ARTIFACT_PATH)" || true
	cd ./test/app && make build
	cp ./test/app/bin/test_app "$(TEST_ARTIFACT_PATH)"
	cp ./test/app/bin/test_app "$(TEST_ARTIFACT_PATH)"

.PHONY: run_test
run_test: clean make_test_app ## run the tests
	$(call print-target)
	PROJECT_PATH=$(PROJECT_PATH) ARTIFACT_PATH=$(TEST_ARTIFACT_PATH) $(GO) test $(BUILD_FLAGS) -v -run $(TEST_REGEX) -p 1 ./test/e2e/

.PHONY: run_test_cover
run_test_cover: clean make_test_app ## run the tests and generate a coverage report
	$(call print-target)
	PROJECT_PATH=$(PROJECT_PATH) ARTIFACT_PATH=$(TEST_ARTIFACT_PATH) $(GO) test $(BUILD_FLAGS) -v -run $(TEST_REGEX) -p 1 -coverprofile=coverage.txt -coverpkg=$(TEST_PACKAGES) ./test/e2e/
	awk $COVERAGE_OMISSION coverage.txt > coverage.out
	rm -f coverage.txt

.PHONY: codecov
codecov: ## process the coverage report and upload it
	$(call print-target)
	codecov -t $(CODECOV_TOKEN) --flags $(CODECOV_FLAG) --file coverage.out

.PHONY: run_test_cover_codecov
run_test_cover_codecov: run_test_cover codecov ## run the tests and process/upload the coverage reports
	$(call print-target)

.PHONY: systemd_install
systemd_install: systemd_uninstall install ## install for systemd-based systems
	$(call print-target)
	cp "$(PROJECT_PATH)/rpm/pmond.service" /usr/lib/systemd/system/
	cp "$(PROJECT_PATH)/rpm/pmond.logrotate" /etc/logrotate.d/pmond
	mkdir -p $(DIR_CONF)
	cp "$(PROJECT_PATH)/config.yml" $(DIR_CONF)
	systemctl enable pmond
	systemctl start pmond
	sh -c "$(PROJECT_PATH)/bin/pmon3 completion bash > /etc/profile.d/pmon3.sh"
	$(PROJECT_PATH)/bin/pmon3 ls
	$(PROJECT_PATH)/bin/pmon3 --help

.PHONY: systemd_uninstall
systemd_uninstall: ## remove from systemd-based systems
	$(call print-target)
	rm -rf $(DIR_CONF) /etc/logrotate.d/pmond /etc/profile.d/pmon3.sh
	systemctl stop pmond || true
	systemctl disable pmond || true

.PHONY: install
install: ## install the binary in the systems executable path
	$(call print-target)
	cp -R bin/pmon* /usr/local/bin/

.PHONY: protogen
protogen: ## generate the protobufs
	$(call print-target)
	protoc protos/*.proto  --go_out=.

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef
