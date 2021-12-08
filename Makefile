.PHONY: default help build docker

SERVICE_NAME=ddns
TARGET_PATH=build
DOCKER_REPO=hub.cn.sang.ink
VERSION=$(shell git describe --tags)

default: help

help: ## show help

	@echo -e "Usage: \n\tmake \033[36m[option]\033[0m"
	@echo -e "Options:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\t\033[36m%-20s\033[0m %s\n", $$1, $$2}'


LD_FLAGS=-ldflags "-X ddns/pkg/config.Version=$(VERSION)"

build: ## build target

	mkdir -p $(TARGET_PATH) $(TARGET_PATH)/bin $(TARGET_PATH)/conf $(TARGET_PATH)/log
	go build $(LD_FLAGS) -o $(TARGET_PATH)/bin/$(SERVICE_NAME) cmd/$(SERVICE_NAME).go
	cp configs/template.yml $(TARGET_PATH)/conf


docker: build ## build docker and push

	docker build -f Dockerfile --no-cache \
		--build-arg DOCKER_PACKAGE_PATH=$(TARGET_PATH) \
		-t $(DOCKER_REPO)/$(SERVICE_NAME):$(VERSION) .
	docker tag $(DOCKER_REPO)/$(SERVICE_NAME):$(VERSION) $(DOCKER_REPO)/$(SERVICE_NAME):latest
	docker push $(DOCKER_REPO)/$(SERVICE_NAME):$(VERSION)
	docker push $(DOCKER_REPO)/$(SERVICE_NAME):latest

run:

	mkdir -p /data/ddns/conf /data/ddns/log
	cp configs/private.yml /data/ddns/conf/ddns.yml
	docker pull $(DOCKER_REPO)/$(SERVICE_NAME):latest
	docker rm -f ddns
	docker run -d \
		--name ddns \
		--restart always \
		--network host \
		-v /data/ddns/conf/ddns.yml:/opt/ddns/conf/ddns.yml \
		-v /data/ddns/log:/opt/ddns/log \
		$(DOCKER_REPO)/$(SERVICE_NAME):latest
