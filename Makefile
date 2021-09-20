VERSION ?= latest
IMAGE ?= ghcr.io/padok-team/git-volume-reloader:$(VERSION)

.DEFAULT_TARGET = help

## help: Display list of commands
.PHONY: help
help: Makefile
	@sed -n 's|^##||p' $< | column -t -s ':' | sed -e 's|^| |'

## build: Build a container image
.PHONY: build
build:
	docker build . -t ${IMAGE}

## push: Push the container image to a registry
.PHONY: push
push:
	docker push ${IMAGE}
