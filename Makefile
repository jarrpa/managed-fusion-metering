.EXPORT_ALL_VARIABLES:
include hack/make-project-vars.mk
include hack/make-tools.mk
include hack/make-bundle-vars.mk


# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.DEFAULT_GOAL := help

all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@./hack/make-help.sh $(MAKEFILE_LIST)

##@ Development

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

lint: ## Run golangci-lint against code.
	docker run --rm -v $(PROJECT_DIR):/app:Z -w /app $(GO_LINT_IMG) golangci-lint run ./...

godeps-update:  ## Run go mod tidy & vendor
	go mod tidy && go mod vendor

test-setup: godeps-update fmt vet ## Run setup targets for tests

go-test: ## Run go test against code.
	./hack/go-test.sh

test: test-setup go-test ## Run go unit tests.

e2e-test: ginkgo ## TODO: Run end to end functional tests.
	@echo "build and run e2e tests"

##@ Build

build: reporter-build ## Build reporter binary

go-build: ## Run go build of reporter against code.
	@GOBIN=${GOBIN} ./hack/go-build.sh

run: fmt vet ## Run reporter program
	go run ./cmd/metering-reporter/main.go

reporter-build: test-setup ## Build reporter image with the application.
	docker build -t ${REPORTER_IMG} -f build/Dockerfile .

reporter-push: ## Push reporter image with the application.
	docker push ${REPORTER_IMG}

