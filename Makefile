.PHONY: test run build
PROJECT_NAME=find-real-activation
VERSION ?= $(VERSION:)

GOPATH ?= $(GOPATH:)
GOFLAGS ?= $(GOFLAGS:)
GO=GOPATH=$(GOPATH) go
export PATH := ${PATH}:${GOPATH}/bin
GO_LINKER_FLAGS ?= -ldflags \

test:
	@echo Running tests
	go test -v ./cmd/...

run:
	@echo Running $(PROJECT_NAME)
	$(GO) run  main.go $(GOFLAGS)


build:
	@echo Building $(PROJECT_NAME)
	$(GO) build

