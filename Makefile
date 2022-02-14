VERSION := $(shell git describe --tags --dirty --always)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GO_LD_FLAGS += -X github.com/jacobtomlinson/containercanairy/internal.Version=$(VERSION)
GO_LD_FLAGS += -X github.com/jacobtomlinson/containercanairy/internal.Buildtime=$(BUILDTIME)
GO_FLAGS = -ldflags "$(GO_LD_FLAGS)"

build:
	go build -o bin/canairy $(GO_FLAGS) .

test:
	go test -v ./...

serve: build
	./bin/canairy

version: build
	./bin/canairy version

help: build
	./bin/canairy help
