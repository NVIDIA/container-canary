VERSION := $(shell git describe --tags --dirty --always)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT := $(shell git rev-parse --short HEAD)

EXECUTABLE=canary

GO_LD_FLAGS += -X github.com/nvidia/container-canary/internal.Version=$(VERSION)
GO_LD_FLAGS += -X github.com/nvidia/container-canary/internal.Buildtime=$(BUILDTIME)
GO_LD_FLAGS += -X github.com/nvidia/container-canary/internal.Commit=$(COMMIT)
GO_FLAGS = -ldflags "$(GO_LD_FLAGS)"

all: test package ## Run tests and build for all platforms

build: ## Build for the current platform
	go build -o bin/$(EXECUTABLE) $(GO_FLAGS) .
	@echo built: bin/$(EXECUTABLE)
	@echo version: $(VERSION)
	@echo commit: $(COMMIT)

package: ## Build for all platforms
	env GOOS=windows GOARCH=amd64 go build -o bin/$(EXECUTABLE)_windows_amd64.exe $(GO_FLAGS) .
	env GOOS=linux GOARCH=amd64 go build -o bin/$(EXECUTABLE)_linux_amd64 $(GO_FLAGS) .
	env GOOS=darwin GOARCH=amd64 go build -o bin/$(EXECUTABLE)_darwin_amd64 $(GO_FLAGS) .
	env GOOS=darwin GOARCH=arm64 go build -o bin/$(EXECUTABLE)_darwin_arm64 $(GO_FLAGS) .
	@echo built:  bin/$(EXECUTABLE)_windows_amd64.exe, bin/$(EXECUTABLE)_linux_amd64, bin/$(EXECUTABLE)_darwin_amd64, bin/$(EXECUTABLE)_darwin_arm64
	@echo version: $(VERSION)
	@echo commit: $(COMMIT)

test: ## Run tests
	go test -v ./...

version:
	@echo version: $(VERSION)

clean: ## Remove previous builds
	rm -f bin/$(EXECUTABLE)*

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
