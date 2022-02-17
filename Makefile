GO                          = $(shell which go)
GOBIN_TOOL                  = $(shell which gobin || echo $(GOPATH)/gobin)
GO_INPUT                    = $(CURDIR)/main.go
GO_OUTPUT                   = $(CURDIR)/bin/$(APP_NAME)
APP_NAME                    ?= helm-repo-updater
GO_TEST_DEFAULT_ARG         = -v ./internal/...

IMAGE_REGISTRY ?= ghcr.io
IMAGE_REPO     ?= docplanner

IMAGE_BUILD_TOOLS = $(IMAGE_REGISTRY)/$(IMAGE_REPO)/helm-repo-updater/build-tools:develop
IMAGE_GIT_REPO_SERVER_TOOL = $(IMAGE_REGISTRY)/$(IMAGE_REPO)/helm-repo-updater/git-repo-server:develop

.PHONY: build
build: clean
	@echo "`date +'%d.%m.%Y %H:%M:%S'` Building $(GO_INPUT)"
	$(GO) build \
    	-ldflags "-s -w" \
    	-o $(GO_OUTPUT) $(GO_INPUT)

.PHONY: clean
clean:
	rm -f $(GO_OUTPUT)

.PHONY: launch-test-deps
launch-test-deps:
ifndef isDevContainer
ifndef isCI
	docker-compose -f test-git-server/docker-compose.yaml up -d
endif
endif
ifdef isDevContainer
	docker-compose -f test-git-server/docker-compose-devcontainer.yaml up -d
endif

.PHONY: clean-test-deps
clean-test-deps:
ifndef isDevContainer
ifndef isCI
	docker-compose -f test-git-server/docker-compose.yaml down
endif
endif
ifdef isDevContainer
	docker-compose -f test-git-server/docker-compose-devcontainer.yaml down
endif
	docker volume prune -f && docker system prune -f

.PHONY: test
test: test-unit

.PHONY: test-benchmark
test-benchmark: launch-test-deps
	$(GO) test ${GO_TEST_DEFAULT_ARG} -cpu 1,2,4,8 -benchmem -bench .

.PHONY: test-unit
test-unit: launch-test-deps
	$(GO) test ${GO_TEST_DEFAULT_ARG}

.PHONY: test-coverage
test-coverage: launch-test-deps
	$(GO) test ${GO_TEST_DEFAULT_ARG} -cover

.PHONY: lint
lint: golangci-lint

.PHONY: golangci-lint
golangci-lint: $(GOBIN_TOOL)
	GOOS=linux $(GOBIN_TOOL) -run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.40.1 run

.PHONY: gofumpt
gofumpt: $(GOBIN_TOOL)
	GOOS=linux $(GOBIN_TOOL) -run mvdan.cc/gofumpt -l -w .

.PHONY: publish-build-tools
publish-build-tools: ## Publish build-tools image
	docker build -f tools/build-tools.Dockerfile -t $(IMAGE_BUILD_TOOLS) .
	docker push $(IMAGE_BUILD_TOOLS)

.PHONY: publish-git-server-tool
publish-git-server-tool: ## Publish git-server-tool image
	docker build -f test-git-server/Dockerfile -t $(IMAGE_GIT_REPO_SERVER_TOOL) .
	docker push $(IMAGE_GIT_REPO_SERVER_TOOL)

$(GOBIN_TOOL):
	go get github.com/myitcv/gobin
	go install github.com/myitcv/gobin
