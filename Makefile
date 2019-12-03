TARGET = ct-template
SOURCE = $(shell find . -type f -name "*.go" -not -name "*_test.go")
VERSION = $(shell cat version.txt)
OUTPUT_DIR = output

.PHONY: all
all: build

.PHONY: setup
setup:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/lint/golint
	curl -o $(OUTPUT_DIR)/ct -sSL https://github.com/coreos/container-linux-config-transpiler/releases/download/v0.9.0/ct-v0.9.0-x86_64-unknown-linux-gnu
	chmod +x $(OUTPUT_DIR)/ct

.PHONY: mod
mod:
	go mod tidy
	go mod vendor

.PHONY: build
build: mod go.mod version.txt $(SOURCE)
	CGO_ENABLED=0 go build -v -o $(OUTPUT_DIR)/$(TARGET) -ldflags "-X main.version=$(VERSION)" .

.PHONY: run
run: build
	$(OUTPUT_DIR)/$(TARGET) ../ct-template/sample | $(OUTPUT_DIR)/ct | jq .

.PHONY: clean
clean:
	-rm $(OUTPUT_DIR)/$(TARGET)

.PHONY: distclean
distclean: clean
	-rm go.sum
	-rm -rf vendor

.PHONY: fmt
fmt:
	goimports -w $$(find . -type d -name 'vendor' -prune -o -type f -name '*.go' -print)

.PHONY: test
test:
	test -z "$$(goimports -l $$(find . -type d -name 'vendor' -prune -o -type f -name '*.go' -print) | tee /dev/stderr)"
	test -z "$$(golint $$(go list ./... | grep -v '/vendor/') | tee /dev/stderr)"
	CGO_ENABLED=0 go test -v ./...
