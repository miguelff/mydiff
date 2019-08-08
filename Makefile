GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
BUILD_DIR=$(shell pwd)/.build
CMD=$(BUILD_DIR)/mydiff

.PHONY: test
test:
	go test -v ./go/...

coverage:
	go test -v ./go/... -coverprofile .coverage
	go tool cover -html .coverage

build: clean
	mkdir -p $(BUILD_DIR)
	$(GO_BUILD_ENV) go build -v -o $(CMD) ./go/cmd

clean:
	rm -rf $(BUILD_DIR)