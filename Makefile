
TEST_PKGS ?= ./...

.PHONY: all
all: test

.PHONY: test
test: export GO_UPGRADE_TEST_RESOURCES := $(PWD)/test
test:
	go test -v $(TEST_PKGS)
