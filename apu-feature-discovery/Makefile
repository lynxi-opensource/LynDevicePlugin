.PHONY: all

IMAGE_BUILD_CMD ?= docker build
IMAGE_PUSH_CMD ?= docker push
CONTAINER_RUN_CMD ?= docker run
BASE_IMAGE_FULL ?= ubuntu:18.04

VERSION ?= 0.1.1
IMAGE_REGISTRY ?= lynxidocker
IMAGE_TAG_NAME ?= $(VERSION)

IMAGE_NAME := apu-feature-discovery
IMAGE_REPO := $(IMAGE_REGISTRY)/$(IMAGE_NAME)
IMAGE_TAG := $(IMAGE_REPO):$(IMAGE_TAG_NAME)

all: build

build: 
	$(IMAGE_BUILD_CMD) --build-arg VERSION=$(VERSION) \
	    --build-arg BASE_IMAGE_FULL=$(BASE_IMAGE_FULL) \
	    -t $(IMAGE_TAG) .

push: build 
	$(IMAGE_PUSH_CMD) $(IMAGE_TAG)

docs:
	godocdown

run:
	# TODO

run_local:
	# TODO