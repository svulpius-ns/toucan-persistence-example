#
# piratetreasure.mk
#
# Make targets which do not call 'drone exec'. These targets are invoked
# by drone steps and can also (optionally) be invoked directly. All
# tooling assumptions must be in place.
#
# WARNING
#
# Any tooling versions (such as protoc) which are not compliant with
# the versions used in the Drone pipeline can lead to unexpected errors.
#
# Proceed with caution.
#

#---------------------
# Build a project binary locally. Feel free to edit as necessary to tailor it
# to the specific needs of your project.
#---------------------
.PHONY: build-piratetreasure
build-piratetreasure:: set-rls-build-info set-ldflags
	@echo ">>-> Building application binaries for $(PROJECT)..."
	go build -ldflags "$(LDFLAGS)" -o dist/$(PROJECT) $(APP_SRC_DIR)/$(PROJECT)

.PHONY: deps-piratetreasure
deps-piratetreasure::
	@echo ">>-> Download dependencies for $(PROJECT)..."
	go mod tidy

.PHONY: lint-piratetreasure
lint-piratetreasure::
	@echo ">>-> Lint checking $(PROJECT)..."
	golangci-lint run -v -c .golangci.yml

.PHONY: coverage-piratetreasure
coverage-piratetreasure::
	@echo ">>-> Running coverage report for $(PROJECT)..."
	nsgocoverreport --summary --configFile=.coverage.yml

.PHONY: ut-piratetreasure
ut-piratetreasure::
	@echo ">>-> Running unit tests for $(PROJECT)..."
	go test -shuffle=on -v -count=1 -coverprofile=coverage.out ./...
	# some tests will only fail with -race enabled; others, without it (!)
	@echo ">>-> Running unit tests with '-race' for $(PROJECT)..."
	CGO_ENABLED=1 go test -shuffle=on -v -count=1 -race ./...

.PHONY: ft-piratetreasure
ft-piratetreasure::
	@echo ">>-> Running functional tests for $(PROJECT)..."
	go test -shuffle=on -cover -v $(shell go list ./functests/...) -functests

.PHONY: benchtest-piratetreasure
benchtest-piratetreasure::
	@echo ">>-> Running benchmark tests for $(PROJECT)..."
	go test -shuffle=on -v -run=XXX --bench . --benchmem ./...

.PHONY: dist-piratetreasure
dist-piratetreasure::
	@echo ">>-> Creating images for $(PROJECT)..."
	@$(eval IMAGE_PATH := $(if $(IMAGE_PATH),$(IMAGE_PATH),artifactory.netskope.io/$(PROJECT)))
	@$(eval TAG := $(if $(TAG),latest,latest-build))
	docker build --build-arg APP_BINARY=dist/piratetreasure -t $(IMAGE_PATH):$(TAG) -f build/piratetreasure/Dockerfile .

.PHONY: clean-piratetreasure
clean-piratetreasure::
	@echo ">>-> Cleaning up $(PROJECT)..."
	rm -rf dist/
	rm -rf docs/sphinx/_build
	rm -rf coverage.out
	rm -rf codeclimate.out
	rm -rf benchmark.out
	rm -rf .drone-env

#----------------------
# Target to clean things beyond just the project directory contents (e.g.,
# docker images).
#----------------------
.PHONY: distclean
distclean:: clean				## Clean build artifacts and deploy images
	@echo ">>-> Checking for images for $(PROJECT)"; \
	deployname='$(PROJECT)\s*latest'; \
	deployimg=$$(docker images | grep -E $${deployname} | awk '{print $$3}'); \
	if [ "$$deployimg" != "" ]; then \
		echo "found image: $$deployimg"; \
		docker rmi -f $$deployimg || { \
			printf "$(WARN) could not remove image '$$deployimg'.\n"; \
		} && { \
			printf "Removed $$deployimg\n"; \
		}; \
	fi; \

#--------------------------
# Set the build info vars when an image is built.
#
# Override these as necessary.
#--------------------------
.PHONY: set-rls-build-info
set-rls-build-info::
	@echo ">>-> Setting release build info..."
	@$(eval BUILD_TIME := $(if $(BUILD_TIME),$(BUILD_TIME),$(shell date +%Y.%m.%d.%H.%M.%S.%N)))
	@$(eval BUILT_BY        := $(if $(BUILT_BY),$(BUILT_BY),$(USER)))
	@$(eval BUILD_HOST      := $(if $(BUILD_HOST),$(BUILD_HOST),$(HOSTNAME)))
	@$(eval PROJECT_VERSION := $(if $(PROJECT_VERSION),$(PROJECT_VERSION),$(APP_VERSION)))
	@$(eval GIT_SHA         := $(if $(GIT_SHA),$(GIT_SHA),$(shell git rev-parse --short=16 HEAD 2>/dev/null || echo no-git-repo)))
	@if [ "$(VERBOSE)" = "yes" ]; then \
		echo "     PROJECT_VERSION = '$(PROJECT_VERSION)'"; \
		echo "     GIT_SHA         = '$(GIT_SHA)'"; \
		echo "     BUILD_HOST      = '$(BUILD_HOST)'"; \
		echo "     BUILD_TIME      = '$(BUILD_TIME)'"; \
		echo "     BUILT_BY        = '$(BUILT_BY)'"; \
	fi

#--------------------------
# This is the central location where the LDFLAGS for the build are set (either
# for host builds or builds done in a container).
#--------------------------
.PHONY: set-ldflags
set-ldflags::
	@echo ">>-> Setting LDFLAGS for build..."
	$(eval LDFLAGS := -X github.com/netskope/piratetreasure/internal/build.AppName=$(PROJECT) \
		-X github.com/netskope/piratetreasure/internal/build.Version=$(PROJECT_VERSION) \
		-X github.com/netskope/piratetreasure/internal/build.GitSha=$(GIT_SHA) \
		-X github.com/netskope/piratetreasure/internal/build.BuildTime=$(BUILD_TIME) \
		-X github.com/netskope/piratetreasure/internal/build.BuiltBy=$(BUILT_BY) \
		-X github.com/netskope/piratetreasure/internal/build.BuildHost=$(BUILD_HOST))
	$(eval LDFLAGS += $(if $(findstring $(strip $(STATIC)),yes),-extldflags '-static',))

#-------------------------
# Target to start metrics.
#-------------------------
.PHONY: start-metrics
start-metrics::
	@command -v docker-compose || (echo "You don't have docker-compose installed; cannot start metrics." 1>&2 && exit 1)
	@mkdir -p build/package/metrics/grafdata
	@mkdir -p build/package/metrics/promdata
	@chmod gou+w build/package/metrics/grafdata
	@chmod gou+w build/package/metrics/promdata
	docker-compose -f build/package/metrics/docker-compose.yml up
