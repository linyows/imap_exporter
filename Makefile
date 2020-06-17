TEST ?= ./...
REVISION = $$(git describe --always)
DATE = $$(LC_ALL=c date -u +%a,\ %d\ %b\ %Y\ %H:%M:%S\ GMT)
LOGLEVEL ?= info
ifeq ("$(shell uname)","Darwin")
NCPU ?= $(shell sysctl hw.ncpu | cut -f2 -d' ')
else
NCPU ?= $(shell cat /proc/cpuinfo | grep processor | wc -l)
endif
TEST_OPTIONS=-timeout 30s -parallel $(NCPU)

default: test

test:
	go test $(TEST) $(TESTARGS) $(TEST_OPTIONS)
	go test -race $(TEST) $(TESTARGS) -coverprofile=coverage.txt -covermode=atomic
lint:
	golint -set_exit_status $(TEST)

.PHONY: default test deps
