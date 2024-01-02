GO ?= go
GOFMT ?= gofmt "-s"
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go")
ROOTDIR=$(shell cd "$(dirname "$0")"; pwd)
WHOAMI=$(shell whoami)
TEST_DIR_LOGS="$(ROOTDIR)/tmp/logs"
TEST_DIR_DATA="$(ROOTDIR)/tmp/data"
TEST_FILE_CONFIG=$(ROOTDIR)/tmp/config-test.yml

all: build

fmt:
	$(GOFMT) -w $(GOFILES)
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;
lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u golang.org/x/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error $(GOFILES)
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w $(GOFILES)
tools:
	@if [ $(GO_VERSION) -gt 15 ]; then \
		$(GO) install golang.org/x/lint/golint@latest; \
		$(GO) install github.com/client9/misspell/cmd/misspell@latest; \
	elif [ $(GO_VERSION) -lt 16 ]; then \
		$(GO) install golang.org/x/lint/golint; \
		$(GO) install github.com/client9/misspell/cmd/misspell; \
	fi
test: build_test
	sudo PMON3_DEBUG=true PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/pmond &
	PMON3_DEBUG=true PMON3_CONF=$(TEST_FILE_CONFIG) ./bin/pmon3 exec bin/test_server
build_test: build
	sudo rm -rf "$(ROOTDIR)/tmp" 
	mkdir -p "$(TEST_DIR_DATA)" "$(TEST_DIR_LOGS)"
	printf '%s\n%s' "data: $(TEST_DIR_DATA)" "logs: $(TEST_DIR_LOGS)" > $(TEST_FILE_CONFIG)
	$(GO) build -o bin/test_server test/test_server.go
build:
	$(GO) mod tidy
	$(GO) build -o bin/pmon3 cmd/pmon3/pmon3.go
	$(GO) build -o bin/pmond cmd/pmond/pmond.go
systemd_install: systemd_uninstall install
	sudo cp "$(ROOTDIR)/rpm/pmond.service" /usr/lib/systemd/system/
	sudo cp "$(ROOTDIR)/rpm/pmond.logrotate" /etc/logrotate.d/pmond
	sudo mkdir -p /var/log/pmond/ /etc/pmon3/config/ /etc/pmon3/data/
	sudo cp "$(ROOTDIR)/config.yml" /etc/pmon3/config/
	sudo systemctl enable pmond
	sudo systemctl start pmond
	sudo sh -c "$(ROOTDIR)/bin/pmon3 completion bash > /etc/profile.d/pmon3.sh"
	$(MAKE) systemd_permissions
	$(ROOTDIR)/bin/pmon3 ls
	$(ROOTDIR)/bin/pmon3 --help
systemd_uninstall: 
	sudo rm -rf /var/log/pmond /etc/pmon3/config /etc/pmon3/data /etc/logrotate.d/pmond /etc/profile.d/pmon3.sh
	sudo systemctl stop pmond
	sudo systemctl disable pmond
systemd_permissions:
	sleep 2
	sudo chown -R root:$(WHOAMI) /etc/pmon3/data/
	sudo chmod 775 /etc/pmon3/data/
	sudo chmod 660 /etc/pmon3/data/*
	sudo chown -R root:$(WHOAMI) /var/log/pmond
	sudo chmod 660 "/var/log/pmond/*" || true
install:
	sudo cp -R bin/pmon* /usr/local/bin/
