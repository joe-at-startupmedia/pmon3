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
test: build run_test

.PHONY: test_cgo
test_cgo: BUILD_FLAGS=$(shell echo '-tags posix_mq,cgo_sqlite')
test_cgo: build_cgo run_test

.PHONY: test_net
test_net: BUILD_FLAGS=$(shell echo '-tags net')
test_net: build run_test

.PHONY: run_test
run_test:
	rm -rf "$(TEST_DIR_DATA)" "$(TEST_DIR_LOGS)"
	mkdir -p "$(TEST_DIR_DATA)" "$(TEST_DIR_LOGS)"
	$(GO) build -o bin/app test/app/app.go
	$(GO) build $(BUILD_FLAGS) -o bin/cli test/cli/cli.go
	cp test/test-apps.config.json "$(TEST_DIR_DATA).."
	cp bin/app "$(TEST_DIR_DATA).."
	$(TEST_VARS) ./bin/pmond > test.log 2>&1 &
	sleep 3 
	@printf "\n\n\033[1mtests that pmond booted from apps config\033[0m\n\n"
	$(TEST_VARS) ./bin/cli ls_assert 2 running
	#
	@printf "\n\n\033[1mtests running additional apps from initial boot\033[0m\n\n"
	$(TEST_VARS) ./bin/cli exec $(ROOTDIR)/bin/app '{"name": "test-server3"}'
	$(TEST_VARS) ./bin/cli ls_assert 3 running
	$(TEST_VARS) ./bin/cli exec $(ROOTDIR)/bin/app '{"name": "test-server4"}'
	$(TEST_VARS) ./bin/cli ls_assert 4 running
	#
	@printf "\n\n\033[1mtests desc command returns nonzero status\033[0m\n\n"
	$(TEST_VARS) ./bin/cli desc 4
	#
	@printf "\n\n\033[1mtests del command removes an application from process list\033[0m\n\n"
	$(TEST_VARS) ./bin/cli del 3 #this is a process id that doesnt exist in the config
	$(TEST_VARS) ./bin/cli ls_assert 3 running
	#
	@printf "\n\n\033[1mtests kill command result in stopping all processes\033[0m\n\n"
	$(TEST_VARS) ./bin/cli kill
	$(TEST_VARS) ./bin/cli ls_assert 3 stopped
	#
	@printf "\n\n\033[1mtests init command restarts all apps\033[0m\n\n"
	$(TEST_VARS) ./bin/cli init all blocking #this will restart all processes (including those speced in the apps config)
	$(TEST_VARS) ./bin/cli ls_assert 3 running
	#
	@printf "\n\n\033[1mtests drop command removes all applications\033[0m\n\n"
	$(TEST_VARS) ./bin/cli drop
	$(TEST_VARS) ./bin/cli ls_assert 0
	#
	@printf "\n\n\033[1mtests init commands boots pmond from app config\033[0m\n\n"
	$(TEST_VARS) ./bin/cli init all blocking
	$(TEST_VARS) ./bin/cli ls_assert 2 running
	#
	@printf "\n\n\033[1mtests that starting and stopping an app works\033[0m\n\n"
	$(TEST_VARS) ./bin/cli drop
	$(TEST_VARS) ./bin/cli exec $(ROOTDIR)/bin/app '{"name": "test-server5"}'
	$(TEST_VARS) ./bin/cli ls_assert 1 running
	$(TEST_VARS) ./bin/cli stop 1
	$(TEST_VARS) ./bin/cli ls_assert 1 stopped
	$(TEST_VARS) ./bin/cli restart 1 '{}'
	$(TEST_VARS) ./bin/cli ls_assert 1 running
	$(TEST_VARS) ./bin/cli drop
	#
	@printf "\n\n\033[1mAll tests passed\033[0m\n\n"
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
