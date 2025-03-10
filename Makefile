ROOT_DIR    = $(shell pwd)
NAMESPACE   = "default"
DEPLOY_NAME = "template-single"
DOCKER_NAME = "template-single"

include ./hack/hack.mk

BINARY=goframechat

.PHONY: run
run:
	@go run main.go

.PHONY: build
build:
	@go build -o $(BINARY)

.PHONY: clean
clean:
	@rm -f $(BINARY)
	@rm -f goframechat.db