GO                          = $(shell which go)
GOBIN_TOOL                  = $(shell which gobin || echo $(GOPATH)/gobin)
GO_INPUT                    = $(CURDIR)/main.go
GO_OUTPUT                   = $(CURDIR)/bin/$(APP_NAME)
APP_NAME                    ?= helm-repo-updater
GO_TEST_DEFAULT_ARG         = -v ./internal/...

IMAGE_REGISTRY	?= ghcr.io
IMAGE_REPO		?= docplanner
VERSION			?= develop


IMAGE						= $(IMAGE_REGISTRY)/$(IMAGE_REPO)/helm-repo-updater:${VERSION}
IMAGE_LATEST				= $(IMAGE_REGISTRY)/$(IMAGE_REPO)/helm-repo-updater:latest
IMAGE_BUILD_TOOLS 			= $(IMAGE_REGISTRY)/$(IMAGE_REPO)/helm-repo-updater/build-tools:${VERSION}
IMAGE_GIT_REPO_SERVER_TOOL 	= $(IMAGE_REGISTRY)/$(IMAGE_REPO)/helm-repo-updater/git-repo-server:${VERSION}

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

.PHONY: docker-build
docker-build: ## Build main image
	docker build -f Dockerfile -t $(IMAGE) -t $(IMAGE_LATEST) .

.PHONY: publish
publish: docker-build ## Publish main image
	docker buildx build --push --platform=linux/amd64,linux/arm64 . -t $(IMAGE) -t $(IMAGE_LATEST)

.PHONY: docker-dev-container
docker-dev-container: ## Build devcontainer image
	docker build -f .devcontainer/Dockerfile .

.PHONY: docker-build-tools
docker-build-tools: ## Build build-tools image
	docker build -f tools/build-tools.Dockerfile -t $(IMAGE_BUILD_TOOLS) .

.PHONY: publish-build-tools
publish-build-tools: docker-build-tools ## Publish build-tools image
	docker push $(IMAGE_BUILD_TOOLS)

.PHONY: docker-git-server-tool
docker-git-server-tool: ## Build git-server-tool image
	docker build -f test-git-server/Dockerfile -t $(IMAGE_GIT_REPO_SERVER_TOOL) .

.PHONY: publish-git-server-tool
publish-git-server-tool: docker-git-server-tool ## Publish git-server-tool image
	docker push $(IMAGE_GIT_REPO_SERVER_TOOL)

$(GOBIN_TOOL):
	go get github.com/myitcv/gobin
	go install github.com/myitcv/gobin
