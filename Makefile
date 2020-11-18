PROJECT := github.com/devfile/integration-tests
ifdef GITCOMMIT
        GITCOMMIT := $(GITCOMMIT)
else
        GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
endif
PKGS := $(shell go list  ./... | grep -v $(PROJECT)/vendor | grep -v $(PROJECT)/tests)
COMMON_LDFLAGS := -X $(PROJECT)/pkg/version.GITCOMMIT=$(GITCOMMIT)
BUILD_FLAGS := -mod=vendor -ldflags="$(COMMON_LDFLAGS)"
CROSS_BUILD_FLAGS := -mod=vendor -ldflags="-s -w $(COMMON_LDFLAGS)"
FILES := integration-tests
TIMEOUT ?= 7200s

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

# Env variable UNIT_TEST_ARGS is used to get control over enabling test flags along with go test.
# For example:
# To enable verbosity export or set env GINKGO_TEST_ARGS like "GINKGO_TEST_ARGS=-v"
UNIT_TEST_ARGS ?=

GINKGO_FLAGS_ALL = $(GINKGO_TEST_ARGS) -randomizeAllSpecs -slowSpecThreshold=$(SLOW_SPEC_THRESHOLD) -timeout $(TIMEOUT)

# Flags for tests that must not be run in parallel.
GINKGO_FLAGS_SERIAL = $(GINKGO_FLAGS_ALL) -nodes=1
# Flags for tests that may be run in parallel
GINKGO_FLAGS=$(GINKGO_FLAGS_ALL) -nodes=$(TEST_EXEC_NODES)


default: bin

.PHONY: bin
bin:
	go build ${BUILD_FLAGS} cmd/odo/odo.go

.PHONY: install
install:
	go install ${BUILD_FLAGS} ./cmd/odo/

# run all validation tests
.PHONY: validate
validate: gofmt check-fit check-vendor vet validate-vendor-licenses sec golint

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
	wwhrd check -q

.PHONY: golint
golint:
	golangci-lint run ./... --timeout 5m

# golint errors are only recommendations
.PHONY: lint
lint:
	golint $(PKGS)

.PHONY: vet
vet:
	go vet $(PKGS)

.PHONY: sec
sec:
	@which gosec 2> /dev/null >&1 || { echo "gosec must be installed to lint code";  exit 1; }
	gosec -severity medium -confidence medium -exclude G304,G204 -quiet  ./...

.PHONY: clean
clean:
	@rm -rf $(FILES)

# install tools used for building, tests and  validations
.PHONY: goget-tools
goget-tools: goget-ginkgo
	mkdir -p $(shell go env GOPATH)/bin
	GOFLAGS='' go get github.com/frapposelli/wwhrd@v0.3.0
	GOFLAGS='' go get github.com/securego/gosec/v2/cmd/gosec@v2.4.0
	# It is not recomended to go get golangci-lint https://github.com/golangci/golangci-lint#go
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.30.0

.PHONY: goget-ginkgo
goget-ginkgo:
	# https://go-review.googlesource.com/c/go/+/198438/
	GOFLAGS='' go get github.com/onsi/ginkgo/ginkgo@v1.14.0

# Run unit tests and collect coverage
.PHONY: test-coverage
test-coverage:
	./scripts/generate-coverage.sh

# compile for multiple platforms
.PHONY: cross
cross:
	./scripts/cross-compile.sh $(CROSS_BUILD_FLAGS)

.PHONY: generate-cli-structure
generate-cli-structure:
	go run cmd/cli-doc/cli-doc.go structure

.PHONY: generate-cli-reference
generate-cli-reference:
	go run cmd/cli-doc/cli-doc.go reference > docs/cli-reference.adoc

.PHONY: configure-installer-tests-cluster
configure-installer-tests-cluster:
	. ./scripts/configure-installer-tests-cluster.sh

# configure cluster to run tests on s390x arch
.PHONY: configure-installer-tests-cluster-s390x
configure-installer-tests-cluster-s390x:
	. ./scripts/configure-installer-tests-cluster-s390x.sh

# configure cluster to run tests on ppc64le arch
.PHONY: configure-installer-tests-cluster-ppc64le
configure-installer-tests-cluster-ppc64le:
	. ./scripts/configure-installer-tests-cluster-ppc64le.sh

.PHONY: test
test:
	go test $(UNIT_TEST_ARGS) -race $(PKGS)

# Run generic integration tests
.PHONY: test-generic
test-generic:
	ginkgo $(GINKGO_FLAGS) -focus="odo generic" tests/integration/

# Run odo login and logout tests
.PHONY: test-cmd-login-logout
test-cmd-login-logout:
	ginkgo $(GINKGO_FLAGS_SERIAL) -focus="odo login and logout command tests" tests/integration/loginlogout/

# Run link and unlink commnad tests against 4.x cluster
.PHONY: test-cmd-link-unlink-4-cluster
test-cmd-link-unlink-4-cluster:
	ginkgo $(GINKGO_FLAGS) -focus="odo link and unlink commnad tests" tests/integration/

# Run odo catalog devfile command tests
.PHONY: test-cmd-devfile-catalog
test-cmd-devfile-catalog:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile catalog command tests" tests/integration/devfile/

# Run odo create devfile command tests
.PHONY: test-cmd-devfile-create
test-cmd-devfile-create:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile create command tests" tests/integration/devfile/

# Run odo push devfile command tests
.PHONY: test-cmd-devfile-push
test-cmd-devfile-push:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile push command tests" tests/integration/devfile/

# Run odo exec devfile command tests
.PHONY: test-cmd-devfile-exec
test-cmd-devfile-exec:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile exec command tests" tests/integration/devfile/

# Run odo exec docker devfile command tests
.PHONY: test-cmd-docker-devfile-exec
test-cmd-docker-devfile-exec:
	ginkgo $(GINKGO_FLAGS) -focus="odo docker devfile exec command tests" tests/integration/devfile/docker

# Run odo status devfile command tests
.PHONY: test-cmd-devfile-status
test-cmd-devfile-status:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile status command tests" tests/integration/devfile/

# Run odo devfile watch command tests
.PHONY: test-cmd-devfile-watch
test-cmd-devfile-watch:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile watch command tests" tests/integration/devfile/

# Run odo devfile app command tests
.PHONY: test-cmd-devfile-app
test-cmd-devfile-app:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile app command tests" tests/integration/devfile/

# Run odo devfile delete command tests
.PHONY: test-cmd-devfile-delete
test-cmd-devfile-delete:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile delete command tests" tests/integration/devfile/

# Run odo devfile registry command tests
.PHONY: test-cmd-devfile-registry
test-cmd-devfile-registry:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile registry command tests" tests/integration/devfile/

# Run odo devfile test command tests
.PHONY: test-cmd-devfile-test
test-cmd-devfile-test:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile test command tests" tests/integration/devfile/
	
# Run odo url devfile command tests
.PHONY: test-cmd-devfile-url
test-cmd-devfile-url:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile url command tests" tests/integration/devfile/

# Run odo debug devfile command tests
.PHONY: test-cmd-devfile-debug
test-cmd-devfile-debug:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile debug command tests" tests/integration/devfile/
	ginkgo $(GINKGO_FLAGS_SERIAL) -focus="odo devfile debug command serial tests" tests/integration/devfile/debug/

# Run odo storage devfile command tests
.PHONY: test-cmd-devfile-storage
test-cmd-devfile-storage:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile storage command tests" tests/integration/devfile/

# Run odo log devfile command tests
.PHONY: test-cmd-devfile-log
test-cmd-devfile-log:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile log command tests" tests/integration/devfile/

# Run odo env devfile command tests
.PHONY: test-cmd-devfile-env
test-cmd-devfile-env:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile env command tests" tests/integration/devfile/

# Run odo config devfile command tests
.PHONY: test-cmd-devfile-config
test-cmd-devfile-config:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile config command tests" tests/integration/devfile/

# Fix https://github.com/openshift/odo/issues/3714 to uncomment docker make target
# # Run odo push docker devfile command tests
# .PHONY: test-cmd-docker-devfile-push
# test-cmd-docker-devfile-push:
# 	ginkgo $(GINKGO_FLAGS) -focus="odo docker devfile push command tests" tests/integration/devfile/docker/

# # Run odo watch docker devfile command tests
# .PHONY: test-cmd-docker-devfile-watch
# test-cmd-docker-devfile-watch:
# 	ginkgo $(GINKGO_FLAGS) -focus="odo docker devfile watch command tests" tests/integration/devfile/docker/

# # Run odo url docker devfile command tests
# .PHONY: test-cmd-docker-devfile-url
# test-cmd-docker-devfile-url:
# 	ginkgo $(GINKGO_FLAGS) -focus="odo docker devfile url command tests" tests/integration/devfile/docker/

# # Run odo docker devfile delete command tests
# .PHONY: test-cmd-docker-devfile-delete
# test-cmd-docker-devfile-delete:
# 	ginkgo $(GINKGO_FLAGS) -focus="odo docker devfile delete command tests" tests/integration/devfile/docker/

# Run odo status docker devfile command tests
#.PHONY: test-cmd-docker-devfile-status
#test-cmd-docker-devfile-status:
#	ginkgo $(GINKGO_FLAGS) -focus="odo docker devfile status command tests" tests/integration/devfile/docker/

# # Run odo catalog devfile command tests
# .PHONY: test-cmd-docker-devfile-catalog
# test-cmd-docker-devfile-catalog:
# 	ginkgo $(GINKGO_FLAGS) -focus="odo docker devfile catalog command tests" tests/integration/devfile/docker/

# # Run odo url docker devfile command tests
# .PHONY: test-cmd-docker-devfile-url-pushtarget
# test-cmd-docker-devfile-url-pushtarget:
# 	ginkgo $(GINKGO_FLAGS) -focus="odo docker devfile url pushtarget command tests" tests/integration/devfile/docker/

# # Run odo test docker devfile command tests
# .PHONY: test-cmd-docker-devfile-test
# test-cmd-docker-devfile-test:
# 	ginkgo $(GINKGO_FLAGS) -focus="odo docker devfile test command tests" tests/integration/devfile/docker/

# Run devfile integration tests
.PHONY: test-integration-devfile
test-integration-devfile:
	ginkgo $(GINKGO_FLAGS) tests/integration/devfile/
	ginkgo $(GINKGO_FLAGS_SERIAL) tests/integration/devfile/debug/

.PHONY: openshiftci-presubmit-unittests
openshiftci-presubmit-unittests:
	./scripts/openshiftci-presubmit-unittests.sh

.PHONY: test-cmd-devfile-describe
test-cmd-devfile-describe:
	ginkgo $(GINKGO_FLAGS) -focus="odo devfile describe command tests" tests/integration/devfile/