PROJECT := github.com/openshift/odo
ifdef GITCOMMIT
        GITCOMMIT := $(GITCOMMIT)
else
        GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
endif

COMMON_GOFLAGS := -mod=vendor
COMMON_LDFLAGS := -X $(PROJECT)/pkg/version.GITCOMMIT=$(GITCOMMIT)
BUILD_FLAGS := $(COMMON_GOFLAGS) -ldflags="$(COMMON_LDFLAGS)"
CROSS_BUILD_FLAGS := $(COMMON_GOFLAGS) -ldflags="-s -w -X $(PROJECT)/pkg/segment.writeKey=R1Z79HadJIrphLoeONZy5uqOjusljSwN $(COMMON_LDFLAGS)"
PKGS := $(shell go list $(COMMON_GOFLAGS)  ./... | grep -v $(PROJECT)/vendor | grep -v $(PROJECT)/tests)
FILES := odo dist
TIMEOUT ?= 14400s

# Env variable TEST_EXEC_NODES is used to pass spec execution type
# (parallel or sequential) for ginkgo tests. To run the specs sequentially use
# TEST_EXEC_NODES=1, otherwise by default the specs are run in parallel on 4 ginkgo test node.
# NOTE: Any TEST_EXEC_NODES value greater than one runs the spec in parallel
# on the same number of ginkgo test nodes.
TEST_EXEC_NODES ?= 2

# Slow spec threshold for ginkgo tests. After this time (in second), ginkgo marks test as slow
SLOW_SPEC_THRESHOLD := 120

# Env variable GINKGO_TEST_ARGS is used to get control over enabling ginkgo test flags against each test target run.
# For example:
# To enable verbosity export or set env GINKGO_TEST_ARGS like "GINKGO_TEST_ARGS=-v"
GINKGO_TEST_ARGS ?=

# ODO_LOG_LEVEL sets the verbose log level for the make tests
export ODO_LOG_LEVEL ?= 4

### description for GINKGO_FLAGS_ALL options
# --noisyPendings=false : suppress printing out a description for each pending spec
# --randomizeAllSpecs : If present, all specs will be permuted. By default Ginkgo will only permute the order of the top level containers.
# --skipMeasurements : If present, Ginkgo will skip any Measure specs you’ve defined.

# GINKGO_FLAGS_ALL = $(GINKGO_TEST_ARGS) --randomizeAllSpecs --skipMeasurements --noisyPendings=false -slowSpecThreshold=$(SLOW_SPEC_THRESHOLD) -timeout $(TIMEOUT)
GINKGO_FLAGS_ALL = $(GINKGO_TEST_ARGS) --noisyPendings=false -slowSpecThreshold=$(SLOW_SPEC_THRESHOLD) -timeout $(TIMEOUT)

# Flags for tests that must not be run in parallel.
GINKGO_FLAGS_SERIAL = $(GINKGO_FLAGS_ALL) -nodes=1
# Flags for tests that may be run in parallel
GINKGO_FLAGS=$(GINKGO_FLAGS_ALL) -nodes=$(TEST_EXEC_NODES)


RUN_GINKGO = GOFLAGS='-mod=vendor' go run $(COMMON_GOFLAGS) github.com/onsi/ginkgo/ginkgo

default: bin

.PHONY: help
help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-26s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: bin
bin: ## build the odo binary
	go build ${BUILD_FLAGS} cmd/odo/odo.go

.PHONY: install
install:
	go install ${BUILD_FLAGS} ./cmd/odo/

.PHONY: validate
validate: gofmt check-fit check-vendor vet validate-vendor-licenses sec golint ## run all validation tests

.PHONY: gofmt
gofmt:
	./scripts/check-gofmt.sh

.PHONY: check-vendor
check-vendor:
	go mod verify

.PHONY: check-fit
check-fit:
	./scripts/check-fit.sh

.PHONY: validate-vendor-licenses
validate-vendor-licenses:
	go run $(COMMON_GOFLAGS) github.com/frapposelli/wwhrd check -q

.PHONY: golint
golint:
	golangci-lint run ./... --timeout 5m

.PHONY: lint
lint: ## golint errors are only recommendations
	golint $(PKGS)

.PHONY: vet
vet:
	go vet $(PKGS)

.PHONY: sec
sec:
	go run $(COMMON_GOFLAGS) github.com/securego/gosec/v2/cmd/gosec -severity medium -confidence medium -exclude G304,G204 -quiet  ./...

.PHONY: clean
clean:
	@rm -rf $(FILES)

.PHONY: goget-tools
goget-tools:
	mkdir -p $(shell go env GOPATH)/bin
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.30.0

### individual devfile command tests
.PHONY: test-cmd-devfile-app
test-cmd-devfile-app: ## Run odo devfile app command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile app command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-catalog
test-cmd-devfile-catalog: ## Run odo catalog devfile command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile catalog command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-config
test-cmd-devfile-config: ## Run odo config devfile command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile config command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-create
test-cmd-devfile-create: ## Run odo create devfile command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile create command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-debug
test-cmd-devfile-debug: ## Run odo debug devfile command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile debug command tests" tests/integration/devfile/
	$(RUN_GINKGO) $(GINKGO_FLAGS_SERIAL) -focus="odo devfile debug command serial tests" tests/integration/devfile/debug/

.PHONY: test-cmd-devfile-delete
test-cmd-devfile-delete: ## Run odo devfile delete command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile delete command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-describe
test-cmd-devfile-describe:
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile describe command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-env
test-cmd-devfile-env: ## Run odo env devfile command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile env command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-exec
test-cmd-devfile-exec: ## Run odo exec devfile command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile exec command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-list
test-cmd-devfile-list: ## Run odo list devfile command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo list with devfile" tests/integration/devfile/

.PHONY: test-cmd-devfile-log
test-cmd-devfile-log: ## Run odo log devfile command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile log command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-push
test-cmd-devfile-push: ## Run odo push devfile command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile push command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-registry
test-cmd-devfile-registry: ## Run odo devfile registry command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile registry command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-status
test-cmd-devfile-status: ## Run odo status devfile command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile status command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-storage
test-cmd-devfile-storage: ## Run odo storage devfile command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile storage command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-test
test-cmd-devfile-test: ## Run odo devfile test command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile test command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-url
test-cmd-devfile-url: ## Run odo url devfile command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile url command tests" tests/integration/devfile/

.PHONY: test-cmd-devfile-watch
test-cmd-devfile-watch: ## Run odo devfile watch command tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile watch command tests" tests/integration/devfile/

### run all devfile integration tests
.PHONY: test-integration-devfile
test-integration-devfile: ## Run devfile integration tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) tests/integration/devfile/
	$(RUN_GINKGO) $(GINKGO_FLAGS_SERIAL) tests/integration/devfile/debug/

### devfile e2e tests
.PHONY: test-e2e-devfile
test-e2e-devfile: ## Run devfile e2e tests: odo devfile supported tests
	$(RUN_GINKGO) $(GINKGO_FLAGS) -focus="odo devfile supported tests" tests/e2escenarios/
