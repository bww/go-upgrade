
export GOPATH := $(GOPATH):$(PWD)

TEST_PKGS := . ./driver/postgres

.PHONY: all test

all: test

test:
	go test -test.v $(TEST_PKGS)
