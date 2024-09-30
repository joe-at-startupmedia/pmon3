GO ?= go
GOFMT ?= gofmt "-s"
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go")
ROOTDIR=$(shell cd "$(dirname "$0")"; pwd)
WHOAMI=$(shell whoami)
TEST_FILE_CONFIG ?= $(ROOTDIR)/test/test-config.yml
TEST_VARS ?= PMON3_CONF=$(TEST_FILE_CONFIG)
TEST_DIR_DATA=$(shell cat $(TEST_FILE_CONFIG) | grep "data_dir:" | cut -d' ' -f2)
TEST_DIR_LOGS=$(shell cat $(TEST_FILE_CONFIG) | grep "logs_dir:" | cut -d' ' -f2)

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
test: build run_e2e_test

.PHONY: test_cgo
test_cgo: BUILD_FLAGS=$(shell echo '-tags posix_mq,cgo_sqlite')
test_cgo: build_cgo run_e2e_test

.PHONY: test_net
test_net: BUILD_FLAGS=$(shell echo '-tags net')
test_net: build run_e2e_test

.PHONY: run_e2e_test
run_e2e_test:
	rm -rf "$(TEST_DIR_DATA)" "$(TEST_DIR_LOGS)"
	mkdir -p "$(TEST_DIR_DATA)" "$(TEST_DIR_LOGS)"
	$(GO) build -o bin/app test/app/app.go
	cp test/test-process.config.json "$(TEST_DIR_DATA).."
	cp bin/app "$(TEST_DIR_DATA).."
	$(TEST_VARS) ./bin/pmond > test.log 2>&1 &
	sleep 3
	$(TEST_VARS) APP_BIN_PATH=$(ROOTDIR) $(GO) test $(BUILD_FLAGS) -v ./test/e2e/
	sleep 3
	@printf "\n\n\033[1mkilling pmond\033[0m\n\n"
	pidof pmond | xargs kill -9

.PHONY: systemd_install
systemd_install: systemd_uninstall install
	cp "$(ROOTDIR)/rpm/pmond.service" /usr/lib/systemd/system/
	cp "$(ROOTDIR)/rpm/pmond.logrotate" /etc/logrotate.d/pmond
	mkdir -p /var/log/pmond/ /etc/pmon3/config/ /etc/pmon3/data/
	cp "$(ROOTDIR)/config.yml" /etc/pmon3/config/
	systemctl enable pmond
	systemctl start pmond
	sh -c "$(ROOTDIR)/bin/pmon3 completion bash > /etc/profile.d/pmon3.sh"
	$(MAKE) systemd_permissions
	$(ROOTDIR)/bin/pmon3 ls
	$(ROOTDIR)/bin/pmon3 --help

.PHONY: systemd_uninstall
systemd_uninstall: 
	rm -rf /var/log/pmond /etc/pmon3/config /etc/pmon3/data /etc/logrotate.d/pmond /etc/profile.d/pmon3.sh
	systemctl stop pmond || true
	systemctl disable pmond || true

.PHONY: systemd_permissions
systemd_permissions:
	sleep 2
	chown -R root:$(WHOAMI) /var/log/pmond
	chmod 660 "/var/log/pmond/*" || true

.PHONY: install
install:
	cp -R bin/pmon* /usr/local/bin/

.PHONY: protogen
protogen:
	protoc pmond/protos/*.proto  --go_out=.
