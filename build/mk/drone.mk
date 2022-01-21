#
# drone.mk
#
# Targets which invoke a drone step or combination of steps via 'drone exec'
#

.ONESHELL:

#--------------------------
# Call the drone cli with dynamic steps. Used by the targets below.
#--------------------------
define DRONE_FUNC =
	$(shell scripts/getnrc)
	$(if $(1), DRONE_STEPS="$(shell echo $(1) | sed 's/[^, ]* */--include=&/g')", )
	cmd="$(DRONE_CMD) $${DRONE_STEPS}"
	pcmd=`echo $${cmd} | sed -e "s#$$password#(password redacted)#g"`
	$(if $(filter "$(VERBOSE)","yes"), @echo "     Command: '$$(echo $${pcmd})'",)
	$(DRONE_CMD) $${DRONE_STEPS}
endef


#--------------------------
# Run particular Drone steps.  ex: `make drone steps=fmt-check,build`
#--------------------------
.PHONY: drone
drone:: env-setup		## Run specific drone pipeline step(s)
	@# run a particular drone step
	[ "$(steps)" ] || {
		echo "$(FAIL) steps not defined" 1>&2
		exit 1
	}
	$(call DRONE_FUNC, $(shell echo $(steps) | tr "," " ")) || {
		echo "$(FAIL) drone exec with steps '$${steps}' failed." 1>&2
		exit 1
	}
	echo "$(OKAY) Hurray! Your drone exec with steps '$${steps}' ran successfully."
	exit 0

#--------------------------
# Run the default Drone CI steps for the event and branch
#--------------------------
.PHONY: ci
ci:: env-setup		## Run the drone pipeline
	@# run the full pipeline for event, and branch
	$(call DRONE_FUNC) || {
		echo "$(FAIL) drone exec failed." 1>&2
		exit 1
	}
	echo "$(OKAY) Hurray! Your drone exec ran successfully."
	exit 0

#----------------------
# Target to clean up the project and bring it back to a state where all build
# artifacts have been removed.
#----------------------
.PHONY: clean
clean::	env-setup		## Clean project build artifacts
	@# clean up the project
	[ "$(VERBOSE)" = "yes" ] && {
		echo ">>-> Cleaning up..."
	}
	$(call DRONE_FUNC, clean-build) || {
		echo -e "$(FAIL) Failure during $@" 1>&2
		exit 1
	}
	echo "$(OKAY) Done cleaning up."
	exit 0

#--------------------------
# Create the application binary
#--------------------------
.PHONY: build-app-binary
build-app-binary: env-setup		## Create app binary
	@# Building the app binary
	[ "$(VERBOSE)" = "yes" ] && {
		echo ">>-> Building the $(PROJECT) application binary..."
	}
	$(call DRONE_FUNC, build) || {
		echo "$(FAIL) Failure during $@" 1>&2
		exit 1
	}
	echo "$(OKAY) Built $(PROJECT)."
	exit 0

#--------------------------
# Create the docker image for deployment
#--------------------------
.PHONY: build-deploy-image
build-deploy-image: env-setup		## Build the deploy image
	@#  Build the deploy image
	[ "$(VERBOSE)" = "yes" ] && {
		echo ">>-> Building $(PROJECT) docker image..."
	}
	$(call DRONE_FUNC, docker-sandbox wait-for-sandbox \
		build build-deploy-image) || {
		echo "$(FAIL) Failure during $@" 1>&2
		exit 1
	}
	echo "$(OKAY) Built $(PROJECT) docker image"
	exit 0

#--------------------------
# Compiles the protobuf
#--------------------------
.PHONY: build-protobuf
build-protobuf:: env-setup		## Builds the protobuf source files into go files
	@# build the protobuf
	[ "$(VERBOSE)" = "yes" ] && {
		echo ">>-> Creating protobuf..."
	}
	$(call DRONE_FUNC, build-protoc-files) || {
		echo "$(FAIL) Could not create protobuf files" 1>&2
		exit 1
	}
	echo "$(OKAY) Created protobuf files."
	exit 0

#--------------------------
# Run untouched
#--------------------------
.PHONY: untouched
untouched:: env-setup		## Run untouched
	@# run untouched
	[ "$(VERBOSE)" = "yes" ] && {
		echo ">>-> Running untouched..."
	}
	$(call DRONE_FUNC, untouched) || {
		echo "$(FAIL) Failure during $@" 1>&2
		exit 1
	}
	echo "$(OKAY) Untouched passed."
	exit 0

#---------------------
# Build the project binaries. The build is done entirely in
# containers; Go does not need to be installed.
#---------------------
.PHONY: project
project:: build-app-binary		## Build project binaries
	@printf "$(OKAY) Successfully built application binaries for $(PROJECT)\n"

#---------------------
# Build the project final deploy image(s). The build is done entirely in
# containers; Go does not need to be installed.
#---------------------
.PHONY: project-images
project-images:: build-deploy-image			## Build all project deploy images
	@printf "$(OKAY) Successfully built final deploy images for $(PROJECT)\n"

#-----------------------------------------------------------------------------
# Test-related targets
#-----------------------------------------------------------------------------

.PHONY: unit
unit: unit-test		## Alias for 'unit-test'

.PHONY: functional
functional:: functional-test		## Alias for 'functional-test'

.PHONY: unit-test
unit-test:: env-setup		## Run project unit tests
	@# test
	[ "$(VERBOSE)" = "yes" ] && {
		echo ">>-> Performing unit tests..."
	}
	$(call DRONE_FUNC, unit-test) || {
		echo "$(FAIL) Unit tests failed." 1>&2
		exit 1
	}
	echo "$(OKAY) Hurray! Your unit tests PASSED."
	exit 0

.PHONY: functional-test
functional-test:: env-setup		## Run project functional tests
	@# test
	[ "$(VERBOSE)" = "yes" ] && {
		echo ">>-> Performing functional tests..."
	}
	$(call DRONE_FUNC, functional-test) || {
		echo "$(FAIL) Functional tests failed." 1>&2
		exit 1
	}
	echo "$(OKAY) Hurray! Your functional tests PASSED."
	exit 0

.PHONY: test-bench
test-bench:: env-setup		## Run project benchmark tests
	@# test benchmarks
	[ "$(VERBOSE)" = "yes" ] && {
		echo ">>-> Performing benchmark tests..."
	}
	$(call DRONE_FUNC, benchmark-test) || {
		echo "$(FAIL) Benchmark tests failed." 1>&2
		exit 1
	}
	echo "$(OKAY) Hurray! Your benchmark tests PASSED."
	exit 0

.PHONY: coverage
coverage:: env-setup		## Run this target to get the coverage percentage in a presentable format
	@# coverage
	[ "$(VERBOSE)" = "yes" ] && {
		echo ">>-> Generating Coverage Report..."
	}
	$(call DRONE_FUNC, unit-test \
		coverage-report) || {
		echo "$(FAIL) Coverage report generation failed." 1>&2
		exit 1
	}
	echo "$(OKAY) Hurray! Your coverage report is published."
	exit 0

.PHONY: lint-check
lint-check:: env-setup		## Perform lint check
	@# verify code layout is correct
	[ "$(VERBOSE)" = "yes" ] && {
		echo ">>-> Performing lint analysis..."
	}
	$(call DRONE_FUNC, lint-check) || {
		echo "$(FAIL) Your code did NOT pass the lint checks." 1>&2
		exit 1
	}
	echo "$(OKAY) Hurray! Your code passed the lint checks."
	exit 0

.PHONY: download-dependencies
download-dependencies:: env-setup		## Download go mod dependencies
	@# Updating the dependencies
	[ "$(VERBOSE)" = "yes" ] && {
		echo ">>-> Performing dependency updates..."
	}
	$(call DRONE_FUNC, download-dependencies) || {
		echo "$(FAIL) Failed download the dependencies." 1>&2
		exit 1
	}
	echo "$(OKAY) Success! All the dependencies are downloaded."
	exit 0

.PHONY: sane-code
sane-code:: env-setup		## Make sure the code is (too) legit (to quit)
	@# run sane-code checks
	[ "$(VERBOSE)" = "yes" ] && {
		echo ">>-> Performing $@ checks..."
	}
	$(call DRONE_FUNC, yaml-lint \
		lint-check \
		unit-test \
		functional-test \
		benchmark-test \
		build \
		coverage-report \
		download-dependencies) || {
		echo "$(FAIL) Failure during $@ checks" 1>&2
		exit 1
	}
	echo "$(OKAY) All basic code sanity checks passed."
	exit 0
