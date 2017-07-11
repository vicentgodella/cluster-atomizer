.DEFAULT_GOAL=test

BUILD_FLAGS=-ldflags=""

# Workaround for GO15VENDOREXPERIMENT bug (https://github.com/golang/go/issues/11659)
ALL_PACKAGES=$(shell go list ./... | grep -v /vendor/ | grep -v /scripts)

build:
	go build -o cluster-atom $(BUILD_FLAGS) github.com/vicentgodella/cluster-atomizer/cmd/atom

install:
	go install $(BUILD_FLAGS) github.com/vicentgodella/cluster-atomizer/cmd/atom

test:
	go test -v $(BUILD_FLAGS) $(ALL_PACKAGES)

cover:
	go test -cover $(BUILD_FLAGS) $(ALL_PACKAGES)

fmt:
	go fmt $(ALL_PACKAGES)
