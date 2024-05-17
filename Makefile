GO ?= go
GOFMT ?= gofmt "-s"
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go")
ROOTDIR=$(shell cd "$(dirname "$0")"; pwd)
WHOAMI=$(shell whoami)
TEST_FILE_CONFIG ?= $(ROOTDIR)/test/test-config.yml
TEST_DIR_DATA=$(shell cat $(TEST_FILE_CONFIG) | grep "data:" | tr -d '"' | cut -d' ' -f2)
TEST_DIR_LOGS=$(shell cat $(TEST_FILE_CONFIG) | grep "logs:" | tr -d '"' | cut -d' ' -f2)

all: build

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt_check
fmt_check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: lint
lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u golang.org/x/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: misspell_check
misspell_check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w $(GOFILES)

.PHONY: tools
tools:
	@if [ $(GO_VERSION) -gt 15 ]; then \
		$(GO) install golang.org/x/lint/golint@latest; \
		$(GO) install github.com/client9/misspell/cmd/misspell@latest; \
	elif [ $(GO_VERSION) -lt 16 ]; then \
		$(GO) install golang.org/x/lint/golint; \
		$(GO) install github.com/client9/misspell/cmd/misspell; \
	fi

.PHONY: test
test: build
	rm -rf "$(TEST_DIR_DATA)" "$(TEST_DIR_LOGS)"
	mkdir -p "$(TEST_DIR_DATA)" "$(TEST_DIR_LOGS)"
	$(GO) build -o bin/app test/app/app.go
	$(GO) build -o bin/cli test/cli/cli.go
	$(TEST_VARS) PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/pmond > test.log 2>&1 &
	sleep 3 
	$(TEST_VARS) PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/cli exec $(ROOTDIR)/bin/app '{"name": "test-server"}'
	$(TEST_VARS) PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/cli ls_assert 1 running
	$(TEST_VARS) PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/cli exec $(ROOTDIR)/bin/app '{"name": "test-server2"}'
	$(TEST_VARS) PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/cli ls_assert 2 running
	$(TEST_VARS) PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/cli desc 2
	$(TEST_VARS) PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/cli del 1
	$(TEST_VARS) PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/cli ls_assert 1 running
	$(TEST_VARS) PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/cli kill
	$(TEST_VARS) PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/cli ls_assert 1 stopped
	$(TEST_VARS) PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/cli init
	$(TEST_VARS) PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/cli ls_assert 1 running
	$(TEST_VARS) PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/cli drop
	$(TEST_VARS) PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/cli ls_assert 0
	pidof pmond | xargs kill -9

.PHONY: build
build:
	$(GO) mod tidy
	CGO_ENABLED=0 $(GO) build -o bin/pmon3 cmd/pmon3/pmon3.go
	CGO_ENABLED=0 $(GO) build -o bin/pmond cmd/pmond/pmond.go

.PHONY: build_cgo
build_cgo:
	$(GO) mod tidy
	CGO_ENABLED=1 $(GO) build -tags $(BUILD_TAGS) -o bin/pmon3 cmd/pmon3/pmon3.go
	CGO_ENABLED=1 $(GO) build -tags $(BUILD_TAGS) -o bin/pmond cmd/pmond/pmond.go

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
