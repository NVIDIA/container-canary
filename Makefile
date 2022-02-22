VERSION := $(shell git describe --tags --dirty --always)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT := $(shell git rev-parse --short HEAD)

GO_LD_FLAGS += -X github.com/jacobtomlinson/containercanary/internal.Version=$(VERSION)
GO_LD_FLAGS += -X github.com/jacobtomlinson/containercanary/internal.Buildtime=$(BUILDTIME)
GO_LD_FLAGS += -X github.com/jacobtomlinson/containercanary/internal.Commit=$(COMMIT)
GO_FLAGS = -ldflags "$(GO_LD_FLAGS)"

build:
	go build -o bin/canary $(GO_FLAGS) .

test:
	go test -v ./...

serve: build
	./bin/canary

version: build
	./bin/canary version

help: build
	./bin/canary help
