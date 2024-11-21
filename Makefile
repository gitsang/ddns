# repo
NAME=ddns

# version
VERSION=$(shell git describe --tags --always --dirty)
BUILD_DATE=$(shell date -u --iso-8601=seconds)

# system
SHELL := /usr/bin/env bash

# go
GO         := go
GOHOSTOS   := $(shell $(GO) env GOHOSTOS)
GOPATH     := $(shell $(GO) env GOPATH)
GO_VERSION := $(shell $(GO) version | awk '{print $$3}')
GOOS       := $(if $(GOOS),$(GOOS),$(shell $(GO) env GOOS))
GOARCH     := $(if $(GOARCH),$(GOARCH),$(shell $(GO) env GOARCH))
GO_MODULE  := $(shell awk '/module/{print $$2}' go.mod)
GIT_COMMIT := $(shell git rev-parse HEAD)
DIST       := dist/$(NAME)/$(GOOS)/$(GOARCH)

# build
LD_FLAGS += -X "main.Version=$(VERSION)"
LD_FLAGS += -X "main.BuildDate=$(BUILD_DATE)"
LD_FLAGS += -X "main.GoVersion=$(GO_VERSION)"
LD_FLAGS += -X "main.GOOS=$(GOOS)"
LD_FLAGS += -X "main.GOARCH=$(GOARCH)"
LD_FLAGS += -X "main.GitCommit=$(shell git rev-parse HEAD)"

# docker
DOCKER = docker
DOCKER_REGISTRY = gitsang


#------------------------------------------------------------------------------#
##@ Debug
#------------------------------------------------------------------------------#


.PHONY: test
## run tests
test:
	$(GO) test -race ./...


.PHONY: run
## run
run:
	$(GO) run ./cmd/ddns \
		$(filter-out $@, $(MAKECMDGOALS))


#------------------------------------------------------------------------------#
##@ Build
#------------------------------------------------------------------------------#


BUILD_DIRS := $(DIST)


$(BUILD_DIRS):
	@mkdir -p $@


.PHONY: build
## build
build:
	$(GO) build -o $(DIST)/$(NAME) ./cmd/ddns


.PHONY: docker
DOCKERFILE ?= Dockerfile
## build docker
docker: $(BUILD_DIRS)
	$(DOCKER) build \
		--no-cache \
		--build-arg VERSION=$(shell git rev-parse HEAD) \
		-t $(DOCKER_REGISTRY)/$(NAME):$(VERSION) \
		-f $(DOCKERFILE) .


.PHONY: publish
## publish docker
publish:
	$(DOCKER) push $(DOCKER_REGISTRY)/$(NAME):$(VERSION)



#------------------------------------------------------------------------------#
##@ Clean
#------------------------------------------------------------------------------#


.PHONY: clean
## clean up git repo
clean:
	git clean -xfd
	rm -fr $(DIST)


#------------------------------------------------------------------------------#
##@ Help
#------------------------------------------------------------------------------#


## display help
help:
	@awk 'BEGIN \
	{ \
		FS = ":.*##"; \
		printf "\nUsage:\n  make \033[36m<target>\033[0m\n" \
	} \
	/^[0-9a-zA-Z\_\-]+:/ \
	{ \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "  \033[36m%-24s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} { lastLine = $$0 } \
	/^##@/ \
	{ \
		printf "\n\033[1m%s\033[0m\n", substr($$0, 5) \
	} ' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
