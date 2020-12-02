GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
DOCKER_BUILD=$(shell pwd)/.docker_build
DOCKER_CMD=$(DOCKER_BUILD)/sso-v2

$(DOCKER_CMD): clean
	mkdir -p $(DOCKER_BUILD)
	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

clean:
	rm -rf $(DOCKER_BUILD)

heroku: $(DOCKER_CMD)
	heroku container:push web

local:
	$(GO_BUILD_ENV) go build -o bin/sso-v2 -v .

PKGS     = $(or $(PKG),$(shell env GO111MODULE=on go list ./...))
TESTPKGS = $(shell env GO111MODULE=on go list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))
TEST_TARGETS := test-default test-bench test-short test-verbose test-race test-cover
.PHONY: $(TEST_TARGETS) test-xml check test tests
test-bench:   ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-short:   ARGS=-short        ## Run only short tests
test-verbose: ARGS=-v            ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race         ## Run tests with race detector
test-cover:   ARGS=-cover        ## Run test with basic coverage
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
test: $(info $(M) running $(NAME:%=% )tests…) ## Run tests
	$(GO_BUILD_ENV) go get github.com/golang/mock/mockgen@v1.4.4
	$(GO_BUILD_ENV) go generate ./...
	go test $(ARGS) $(TESTPKGS)