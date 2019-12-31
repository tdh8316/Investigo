# docker variables
IMAGE 			:= x0rzkov/investigo
TAG 			:= standalone
TAG_CHROMIUM 	:= chromium

# build variables
VCS_REF 		:= $(shell git describe HEAD --always)
BRANCH 			:= $(shell git rev-parse --abbrev-ref HEAD | tr / -)
BUILD_DATE		 = $(shell TZ=UTC date +%Y-%m-%dT%H:%M:%SZ)

## test			:	test.
test:
	true

## run			:	run generator (requires golang to be already installed).
.PHONY: run
run: deps
	@go run --race *.go

## build			:	build generator (requires golang to be already installed).
.PHONY: build
build: deps
	@go build -v

## deps			:	install dependencies.
.PHONY: deps
deps:
	@go mod tidy

## docker-image		:	build the standalone image and tag to latest.
.PHONY: docker-image
docker-image:
	@docker build --build-arg BUILD_DATE=$(BUILD_DATE) --build-arg VCS_REF=$(VCS_REF) -t "$(IMAGE):$(TAG)-$(VCS_REF)" -f Dockerfile .
	@docker tag $(IMAGE):$(TAG)-$(VCS_REF) $(IMAGE):$(TAG)-latest

## docker-chromium	:	build image with chromium and tag to latest.
.PHONY: docker-chromium
docker-chromium:
	@docker build --build-arg BUILD_DATE=$(BUILD_DATE) --build-arg VCS_REF=$(VCS_REF) -t "$(IMAGE):$(TAG_CHROMIUM)-$(VCS_REF)" -f Dockerfile.chromium .
	@docker tag $(IMAGE):$(TAG_CHROMIUM)-$(VCS_REF) $(IMAGE):$(TAG_CHROMIUM)-latest

# docker-run : run investigo container (by default, it displays help commans)
.PHONY: docker-run
docker-run:
	@docker run -ti -v $(PWD):/opt/twint-docker/data "$(IMAGE):$(VCS_REF)"

## docker-push		:	push docker image.
.PHONY: docker-push
docker-push:
	@docker push $(IMAGE):$(VCS_REF)
	@docker push $(IMAGE):latest

## compose-build		:	build with docker-compose.
.PHONY: compose-build
compose-build:
	@docker-compose build --build-arg BUILD_DATE=$(BUILD_DATE) --build-arg VCS_REF=$(VCS_REF)

## help			:	Print commands help.
.PHONY: help
help : Makefile
	@sed -n 's/^##//p' $<

# https://stackoverflow.com/a/6273809/1826109
%:
	@:
