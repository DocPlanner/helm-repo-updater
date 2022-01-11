GO                          = $(shell which go)
GOBIN_TOOL                  = $(shell which gobin || echo $(GOPATH)/gobin)
GO_INPUT                    = $(CURDIR)/main.go
GO_OUTPUT                   = $(CURDIR)/bin/$(APP_NAME)
APP_NAME                    =? helm-repo-updater
GO_TEST_DEFAULT_ARG         = -v ./internal/...

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
	docker-compose -f test-git-server/docker-compose.yaml up -d --build

.PHONY: clean-test-deps
clean-test-deps:
	docker-compose -f test-git-server/docker-compose.yaml down
	docker volume prune -f && docker system prune -f

.PHONY: test
test: test-unit

.PHONY: test-benchmark
test-benchmark:
	$(GO) test ${GO_TEST_DEFAULT_ARG} -cpu 1,2,4,8 -benchmem -bench .

.PHONY: test-unit
test-unit:
	$(GO) test ${GO_TEST_DEFAULT_ARG}

.PHONY: test-coverage
test-coverage:
	$(GO) test ${GO_TEST_DEFAULT_ARG} -cover

.PHONY: lint
lint: golangci-lint

.PHONY: golangci-lint
golangci-lint: $(GOBIN_TOOL)
	GOOS=linux $(GOBIN_TOOL) -run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.40.1 run

.PHONY: gofumpt
gofumpt: $(GOBIN_TOOL)
	GOOS=linux $(GOBIN_TOOL) -run mvdan.cc/gofumpt -l -w .

$(GOBIN_TOOL):
	go get github.com/myitcv/gobin
	go install github.com/myitcv/gobin