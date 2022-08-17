SPIS_TAG_NAME ?= next
SPIS_IMAGE_TAG_BASE ?= quay.io/redhat-appstudio/service-provider-integration-scm-file-retriever-server
SPIS_IMG ?= $(SPIS_IMAGE_TAG_BASE):$(SPIS_TAG_NAME)

SHELL := bash
.SHELLFLAGS = -ec
.ONESHELL:
.DEFAULT_GOAL := help
ifndef VERBOSE
  MAKEFLAGS += --silent
endif

ENVTEST_K8S_VERSION = 1.21
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

##@ General

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

test: fmt fmt_license vet envtest ## Run the unit tests
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test ./... -cover

run: ## Run the binary
	cd server
	go run server.go

vet: fmt fmt_license ## Run go vet against code.
	go vet ./...

##@ Build

build: fmt fmt_license vet ## Builds the binary
	cd server
	go build -o bin/service-provider-integration-scm-file-retriever-server server.go

docker-build: fmt fmt_license vet ## Builds the docker image. Use the SPI_IMG env var to override the image tag
	cd server
	docker build -t ${SPIS_IMG} .

docker-push: docker-build ## Pushes the image. Use the SPI_IMG env var to override the image tag
	docker push ${SPIS_IMG}

fmt:
  ifneq ($(shell command -v goimports 2> /dev/null),)
	  find . -not -path '*/\.*' -name '*.go' -exec goimports -w {} \;
  else
	  @echo "WARN: goimports is not installed -- formatting using go fmt instead."
	  @echo "      Please install goimports to ensure file imports are consistent."
	  go fmt -x ./...
  endif

fmt_license:
  ifneq ($(shell command -v addlicense 2> /dev/null),)
	  @echo 'addlicense -v -f license_header.txt **/*.go'
	  addlicense -v -f license_header.txt $$(find . -not -path '*/\.*' -name '*.go')
  else
	  $(error addlicense must be installed for this rule: go install github.com/google/addlicense@latest)
  endif

### check_fmt: Checks the formatting on files in repo
check_fmt:
  ifeq ($(shell command -v goimports 2> /dev/null),)
	  $(error "goimports must be installed for this rule" && exit 1)
  endif
  ifeq ($(shell command -v addlicense 2> /dev/null),)
	  $(error "error addlicense must be installed for this rule: go install github.com/google/addlicense@latest")
  endif
	  if [[ $$(find . -not -path '*/\.*' -name '*.go' -exec goimports -l {} \;) != "" ]]; then \
	    echo "Files not formatted; run 'make fmt'"; exit 1 ;\
	  fi ;\
	  if ! addlicense -check -f license_header.txt $$(find . -not -path '*/\.*' -name '*.go'); then \
	    echo "Licenses are not formatted; run 'make fmt_license'"; exit 1 ;\
	  fi

lint: ## Run the linter on the codebase
  ifeq ($(shell command -v golangci-lint 2> /dev/null),)
	  $(error "golangci-lint must be installed for this rule" && exit 1)
  endif
	golangci-lint run

check: check_fmt lint test ## Check that the code conforms to all requirements for commit. Formatting, licenses, vet, tests and linters

ENVTEST = $(shell pwd)/bin/setup-envtest
envtest: ## Download envtest-setup locally if necessary.
	$(call go-get-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest@latest)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
