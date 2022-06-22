.PHONY: default help build docker

SERVICE_NAME=ddns
TARGET_PATH=build
DOCKER_REPO=gitsang
VERSION=$(shell git describe --tags)

default: help

help: ## show help

	@echo -e "Usage: \n\tmake \033[36m[option]\033[0m"
	@echo -e "Options:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\t\033[36m%-20s\033[0m %s\n", $$1, $$2}'


LD_FLAGS=-ldflags "-X ddns/pkg/config.Version=$(VERSION)"

clean: ## clean build target

	rm -fr $(TARGET_PATH)
	rm -f *.tar.gz

build: clean ## build target

	mkdir -p $(TARGET_PATH) $(TARGET_PATH)/bin $(TARGET_PATH)/conf $(TARGET_PATH)/log
	go build $(LD_FLAGS) -o $(TARGET_PATH)/bin/$(SERVICE_NAME) cmd/$(SERVICE_NAME).go
	cp configs/template.yml $(TARGET_PATH)/conf
	cp configs/ddns.service $(TARGET_PATH)/conf

tag: build ## make tgz package

	cp -r $(TARGET_PATH) $(SERVICE_NAME)
	tar zcvf $(SERVICE_NAME).$(VERSION).tar.gz $(SERVICE_NAME)
	rm -fr $(SERVICE_NAME)

docker: build ## build docker

	docker build -f Dockerfile --no-cache \
		--build-arg DOCKER_PACKAGE_PATH=$(TARGET_PATH) \
		-t $(DOCKER_REPO)/$(SERVICE_NAME):$(VERSION) .
	docker tag $(DOCKER_REPO)/$(SERVICE_NAME):$(VERSION) $(DOCKER_REPO)/$(SERVICE_NAME):latest

publish: docker ## publish docker to docker repo

	docker push $(DOCKER_REPO)/$(SERVICE_NAME):$(VERSION)
	docker push $(DOCKER_REPO)/$(SERVICE_NAME):latest

pull: ## pull latest published docker

	docker pull $(DOCKER_REPO)/$(SERVICE_NAME):latest

install: ## install by systemd

	cp $(TARGET_PATH)/bin/$(SERVICE_NAME) /usr/local/bin/
	cp configs/ddns.service /etc/systemd/system/
	cp configs/template.yml /usr/local/etc/ddns/ddns.yml
	mkdir -p /var/log/ddns/

docker-install: docker ## build docker and run

	mkdir -p /data/ddns/conf /data/ddns/log
	docker rm -f ddns
	docker run -d \
		--name ddns \
		--restart always \
		--network host \
		-v /data/ddns/conf/ddns.yml:/opt/ddns/conf/ddns.yml \
		-v /data/ddns/log:/opt/ddns/log \
		$(DOCKER_REPO)/$(SERVICE_NAME):latest
