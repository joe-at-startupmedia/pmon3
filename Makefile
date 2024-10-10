GO ?= go
GOFMT ?= gofmt "-s"
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go")
PROJECT_PATH=$(shell cd "$(dirname "$0")"; pwd)
DIR_CONF ?= /etc/pmon3/config
WHOAMI=$(shell whoami)
TEST_REGEX := $(or $(TEST_REGEX),"Test")
TEST_FILE_CONFIG ?= $(PROJECT_PATH)/test/e2e/config/test-config.core.yml
TEST_DIR_DATA=$(shell cat $(TEST_FILE_CONFIG) | grep "^  data:" | cut -d' ' -f4)
TEST_DIR_LOGS=$(shell cat $(TEST_FILE_CONFIG) | grep "^  logs:" | cut -d' ' -f4)
TEST_ARTIFACT_PATH=$(shell dirname "$(TEST_DIR_DATA)")
DEFAULT_TEST_PACKAGES := "pmon3/cli/...,pmon3/conf,pmon3/pmond/controller/...,pmon3/pmond/db,pmon3/pmond/god,pmon3/pmond/model,pmon3/pmond/observer,pmon3/pmond/process,pmon3/pmond/repo,pmon3/pmond/shell"
TEST_PACKAGES := $(or $(TEST_PACKAGES),$(DEFAULT_TEST_PACKAGES))

all: build

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt_check
fmt_check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		@printf "Please run 'make fmt' and commit the result:"; \
		@printf "$${diff}"; \
		exit 1; \
	fi;

.PHONY: lint
lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) install golang.org/x/lint/golint@latest
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: misspell_check
misspell_check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) install github.com/client9/misspell/cmd/misspell@latest; \
	fi
	misspell -error $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) install github.com/client9/misspell/cmd/misspell@latest; \
	fi
	misspell -w $(GOFILES)

.PHONY: betteralign_check
betteralign_check:
	@hash betteralign > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) install github.com/dkorunic/betteralign/cmd/betteralign@v0.5.1; \
	fi
	betteralign ./...

.PHONY: betteralign
betteralign:
	@hash betteralign > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) install github.com/dkorunic/betteralign/cmd/betteralign@v0.5.1; \
	fi
	betteralign -apply ./...

.PHONY: tools
tools:
	$(GO) install golang.org/x/lint/golint@latest
	$(GO) install github.com/client9/misspell/cmd/misspell@latest
	$(GO) install github.com/dkorunic/betteralign/cmd/betteralign@v0.5.1

.PHONY: base_build
base_build: misspell_check betteralign_check
	$(GO) mod tidy
	$(ENV_VARS) $(GO) build $(BUILD_FLAGS) -o bin/pmon3 cmd/pmon3/pmon3.go
	$(ENV_VARS) $(GO) build $(BUILD_FLAGS) -o bin/pmond cmd/pmond/pmond.go

.PHONY: build
build: ENV_VARS=CGO_ENABLED=0
build: base_build

.PHONY: build_cgo
build_cgo: ENV_VARS=CGO_ENABLED=1
build_cgo: base_build

.PHONY: test
test: build run_test

.PHONY: test_cgo
test_cgo: BUILD_FLAGS=$(shell echo '-tags posix_mq,cgo_sqlite')
test_cgo: build_cgo run_test

.PHONY: test_net
test_net: BUILD_FLAGS=$(shell echo '-tags net')
test_net: build run_test

.PHONY: run_test
run_test:
	rm -rf "$(TEST_ARTIFACT_PATH)"
	mkdir -p "$(TEST_ARTIFACT_PATH)"
	cd ./test/app && make build
	cp ./test/app/bin/test_app "$(TEST_ARTIFACT_PATH)"
	PROJECT_PATH=$(PROJECT_PATH) ARTIFACT_PATH=$(TEST_ARTIFACT_PATH) $(GO) test $(BUILD_FLAGS) -v -run $(TEST_REGEX) -p 1 -coverprofile=coverage.txt -coverpkg=$(TEST_PACKAGES) ./test/e2e/

.PHONY: systemd_install
systemd_install: systemd_uninstall install
	cp "$(PROJECT_PATH)/rpm/pmond.service" /usr/lib/systemd/system/
	cp "$(PROJECT_PATH)/rpm/pmond.logrotate" /etc/logrotate.d/pmond
	mkdir -p $(DIR_CONF)
	cp "$(PROJECT_PATH)/config.yml" $(DIR_CONF)
	systemctl enable pmond
	systemctl start pmond
	sh -c "$(PROJECT_PATH)/bin/pmon3 completion bash > /etc/profile.d/pmon3.sh"
	$(MAKE) systemd_permissions
	$(PROJECT_PATH)/bin/pmon3 ls
	$(PROJECT_PATH)/bin/pmon3 --help

.PHONY: systemd_uninstall
systemd_uninstall: 
	rm -rf $(DIR_CONF) /etc/logrotate.d/pmond /etc/profile.d/pmon3.sh
	systemctl stop pmond || true
	systemctl disable pmond || true

.PHONY: systemd_permissions
systemd_permissions:
	sleep 2
	chown -R root:$(WHOAMI) $(DIR_LOGS)
	chmod 660 "$(DIR_LOGS)/*" || true

.PHONY: install
install:
	cp -R bin/pmon* /usr/local/bin/

.PHONY: protogen
protogen:
	protoc pmond/protos/*.proto  --go_out=.
