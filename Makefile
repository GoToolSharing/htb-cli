BINARY_NAME=htb-cli

PLATFORMS=darwin linux windows

ARCHITECTURES=amd64 arm64

OUTPUT_DIR=build

all: test coverage $(PLATFORMS)

$(PLATFORMS):
	for arch in $(ARCHITECTURES); do \
		GOOS=$@ GOARCH=$$arch go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-$@-$$arch; \
	done

build:
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME)
	
test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf $(OUTPUT_DIR) coverage.out coverage.html

.PHONY: all test coverage clean $(PLATFORMS)
