# ====================================================================================
# Setup Project

PROJECT_NAME := x-metrics
PROJECT_REPO := github.com/crossplane-contrib/$(PROJECT_NAME)

PLATFORMS ?= linux_amd64 linux_arm64
# -include will silently skip missing files, which allows us
# to load those files with a target in the Makefile. If only
# "include" was used, the make command would fail and refuse
# to run a target until the include commands succeeded.
-include build/makelib/common.mk

# ====================================================================================
# Setup Output

S3_BUCKET ?= crossplane.releases
-include build/makelib/output.mk

# ====================================================================================
# Setup Go

# Set a sane default so that the nprocs calculation below is less noisy on the initial
# loading of this file
NPROCS ?= 1

# each of our test suites starts a kube-apiserver and running many test suites in
# parallel can lead to high CPU utilization. by default we reduce the parallelism
# to half the number of CPU cores.
GO_TEST_PARALLEL := $(shell echo $$(( $(NPROCS) / 2 )))

GO_STATIC_PACKAGES = $(GO_PROJECT)/cmd/x-metrics 
# GO_TEST_PACKAGES = $(GO_PROJECT)/test/e2e
GO_LDFLAGS += -X $(GO_PROJECT)/internal/version.version=$(VERSION)
GO_SUBDIRS += api
GO111MODULE = on
GOLANGCILINT_VERSION = 1.53.3
-include build/makelib/golang.mk

# ====================================================================================
# Setup Kubernetes tools

KIND_VERSION = v0.20.0
UP_VERSION = v0.18.0
UP_CHANNEL = stable
-include build/makelib/k8s_tools.mk

# ====================================================================================
# Setup Images
# Due to the way that the shared build logic works, images should
# all be in folders at the same level (no additional levels of nesting).

REGISTRY_ORGS = xpkg.upbound.io/crossplane-contrib
IMAGES = x-metrics
-include build/makelib/imagelight.mk

# ====================================================================================
# ControllerGen
# Due to the way that the shared build logic works, images should
# all be in folders at the same level (no additional levels of nesting).

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
CONTROLLER_TOOLS_VERSION ?= v0.9.2
CRD_DIR = package/crds

# ====================================================================================
# Targets

# run `make help` to see the targets and options

# We want submodules to be set up the first time `make` is run.
# We manage the build/ folder and its Makefiles as a submodule.
# The first time `make` is run, the includes of build/*.mk files will
# all fail, and this target will be run. The next time, the default as defined
# by the includes will be run instead.

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

fallthrough: submodules controller-gen
	@echo Initial setup complete. Running make again . . .
	@make

manifests:
	@$(WARN) Deprecated. Please run make generate instead.

crds.clean:
	@$(INFO) cleaning generated CRDs
	@find $(CRD_DIR) -name '*.yaml' -exec sed -i.sed -e '1,1d' {} \; || $(FAIL)
	@find $(CRD_DIR) -name '*.yaml.sed' -delete || $(FAIL)
	@$(OK) cleaned generated CRDs

generate.init: gen-generate gen-crds
generate.run:  gen-kustomize-crds gen-chart-license

gen-chart-license:
	@cp -f LICENSE cluster/charts/x-metrics/LICENSE

generate.done: crds.clean

## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
gen-generate:
	@$(INFO) generate code
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."
	@$(OK) generate code

gen-crds:
	@$(INFO) generate CRDs
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd paths="./..." output:crd:artifacts:config=package/crds
	@cp package/crds/* cluster/charts/x-metrics/templates/crds/
	@cp package/crds/* cluster/crds/
	@$(OK) generate CRDs

gen-kustomize-crds:
	@$(INFO) Adding all CRDs to Kustomize file for local development
	@rm -rf cluster/kustomization.yaml
	@echo "# This kustomization can be used to remotely install all Crossplane CRDs" >> cluster/kustomization.yaml
	@echo "# by running kubectl apply -k https://github.com/crossplane-contrib/x-metrics//cluster?ref=master" >> cluster/kustomization.yaml
	@echo "resources:" >> cluster/kustomization.yaml
	@find $(CRD_DIR) -type f -name '*.yaml' | sort | \
		while read filename ;\
		do echo "- $${filename#*/}" >> cluster/kustomization.yaml \
		; done
	@$(OK) All CRDs added to Kustomize file for local development

# Update the submodules, such as the common build scripts.
submodules: controller-gen
	@git submodule sync
	@git submodule update --init --recursive

# Install CRDs into a cluster. This is for convenience.
install-crds: $(KUBECTL) reviewable
	$(KUBECTL) apply -f $(CRD_DIR)

# Uninstall CRDs from a cluster. This is for convenience.
uninstall-crds:
	$(KUBECTL) delete -f $(CRD_DIR)

# NOTE(hasheddan): the build submodule currently overrides XDG_CACHE_HOME in
# order to force the Helm 3 to use the .work/helm directory. This causes Go on
# Linux machines to use that directory as the build cache as well. We should
# adjust this behavior in the build submodule because it is also causing Linux
# users to duplicate their build cache, but for now we just make it easier to
# identify its location in CI so that we cache between builds.
go.cachedir:
	@go env GOCACHE

# This is for running out-of-cluster locally, and is for convenience. Running
# this make target will print out the command which was used. For more control,
# try running the binary directly with different arguments.
run: go.build
	@$(INFO) Running x-metrics locally out-of-cluster . . .
	@# To see other arguments that can be provided, run the command with --help instead
	$(GO_OUT_DIR)/$(PROJECT_NAME)

.PHONY: manifests submodules fallthrough test-integration run install-crds uninstall-crds gen-kustomize-crds e2e-tests-compile e2e.test.images

# ====================================================================================
# Special Targets

define CROSSPLANE_MAKE_HELP
Crossplane Targets:
    submodules         Update the submodules, such as the common build scripts.
    run                Run crossplane locally, out-of-cluster. Useful for development.

endef
# The reason CROSSPLANE_MAKE_HELP is used instead of CROSSPLANE_HELP is because the crossplane
# binary will try to use CROSSPLANE_HELP if it is set, and this is for something different.
export CROSSPLANE_MAKE_HELP

crossplane.help:
	@echo "$$CROSSPLANE_MAKE_HELP"

help-special: crossplane.help

.PHONY: crossplane.help help-special
