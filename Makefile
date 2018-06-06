.DEFAULT_GOAL=test
UNAME=$(shell uname)
PREFIX=github.com/gangachris/hlf-cli
GOPATH=$(shell go env GOPATH)
GOVERSION=$(shell go version)
BUILD_SHA=$(shell git rev-parse HEAD)
LLDB_SERVER=$(shell which lldb-server)

ifeq "$(UNAME)" "Darwin"
    BUILD_FLAGS=-ldflags="-s -X main.Build=$(BUILD_SHA)"
else
    BUILD_FLAGS=-ldflags="-X main.Build=$(BUILD_SHA)"
endif

# Workaround for GO15VENDOREXPERIMENT bug (https://github.com/golang/go/issues/11659)
ALL_PACKAGES=$(shell go list ./... | grep -v /vendor/ | grep -v /scripts)

# We must compile with -ldflags="-s" to omit
# DWARF info on OSX when compiling with the
# 1.5 toolchain. Otherwise the resulting binary
# will be malformed once we codesign it and
# unable to execute.
# See https://github.com/golang/go/issues/11887#issuecomment-126117692.
ifeq "$(UNAME)" "Darwin"
	TEST_FLAGS=-count 1 -exec=$(shell pwd)/scripts/testsign
	export PROCTEST=lldb
	DARWIN="true"
else
	TEST_FLAGS=-count 1
endif


install:
	go install $(BUILD_FLAGS) github.com/gangachris/hlf
