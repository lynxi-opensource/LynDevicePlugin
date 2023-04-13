.PHONY: all

IMAGE_BUILD_CMD ?= docker build
IMAGE_PUSH_CMD ?= docker push
CONTAINER_RUN_CMD ?= docker run
BASE_IMAGE_FULL ?= ubuntu:18.04

VERSION ?= 1.4.0
IMAGE_REGISTRY ?= lynxidocker
IMAGE_TAG_NAME ?= $(VERSION)

IMAGE_NAME := apu-feature-discovery
IMAGE_REPO := $(IMAGE_REGISTRY)/$(IMAGE_NAME)
IMAGE_TAG := $(IMAGE_REPO):$(IMAGE_TAG_NAME)

build-amd64: 
	$(IMAGE_BUILD_CMD) --build-arg VERSION=$(VERSION) \
	    --build-arg BASE_IMAGE_FULL=$(BASE_IMAGE_FULL) \
	    -t $(IMAGE_TAG)-amd64 .

push-amd64: build-amd64
	$(IMAGE_PUSH_CMD) $(IMAGE_TAG)-amd64

build-arm64: 
	$(IMAGE_BUILD_CMD) --build-arg VERSION=$(VERSION) \
	    --build-arg BASE_IMAGE_FULL=$(BASE_IMAGE_FULL) \
	    -t $(IMAGE_TAG)-arm64 .

push-arm64: build-arm64
	$(IMAGE_PUSH_CMD) $(IMAGE_TAG)-arm64

docker-manifest:
	docker manifest rm $(IMAGE_TAG)
	docker manifest create $(IMAGE_TAG) $(IMAGE_TAG)-amd64 $(IMAGE_TAG)-arm64
	docker manifest annotate $(IMAGE_TAG) $(IMAGE_TAG)-amd64 --os linux --arch amd64
	docker manifest annotate $(IMAGE_TAG) $(IMAGE_TAG)-arm64 --os linux --arch arm64
	docker manifest push $(IMAGE_TAG)

docs:
	godocdown

run:
	# TODO

run_local:
	# TODO