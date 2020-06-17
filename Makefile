TEST ?= .
ifeq ("$(shell uname)","Darwin")
NCPU ?= $(shell sysctl hw.ncpu | cut -f2 -d' ')
else
NCPU ?= $(shell cat /proc/cpuinfo | grep processor | wc -l)
endif
TEST_OPTIONS=-timeout 30s -parallel $(NCPU)
PREFIX  ?= $(shell pwd)
BIN_DIR ?= $(shell pwd)

default: test

test:
	go test $(TEST) $(TESTARGS) $(TEST_OPTIONS)
	go test -race $(TEST) $(TESTARGS) -coverprofile=coverage.txt -covermode=atomic
lint:
	golint -set_exit_status $(TEST)

cross_build: promu
	promu crossbuild

cross_tarball: cross_build
	@rm -rf $(PREFIX)/.tarballs
	promu crossbuild tarballs

release: cross_tarball
	promu checksum .tarballs
	promu release .tarballs

build: promu
	promu build --prefix $(PREFIX)

tarball: promu
	promu tarball --prefix $(PREFIX) $(BIN_DIR)

promu:
	go get -u github.com/prometheus/promu

.PHONY: default test deps
