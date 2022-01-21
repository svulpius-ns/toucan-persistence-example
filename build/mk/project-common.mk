#
# project-common.mk
#
# Common includes for building Go binaries.
#

#---------------------------------------------------
# Default targets used elsewhere in goproject. By putting the defaults here,
# projects can have them wired into their top-level Makefile but they do
# nothing. When they layer on something new to their project (e.g., protobuf
# usage), nothing needs to change. Add others as necessary.
#---------------------------------------------------

# Color formatting of output
COLOR_FG    := \033[0m
BLACK_ON    := \033[0;30m
RED_ON      := \033[0;31m
GREEN_ON    := \033[0;32m
ORANGE_ON   := \033[0;33m
BLUE_ON     := \033[0;34m
PURPLE_ON   := \033[0;35m
CYAN_ON     := \033[0;36m
LTGRAY_ON   := \033[0;37m
DKGRAY_ON   := \033[0;30m
LTRED_ON    := \033[0;31m
LTGREEN_ON  := \033[0;32m
YELLOW_ON   := \033[0;33m
LTBLUE_ON   := \033[0;34m
LTPURPLE_ON := \033[0;35m
LTCYAN_ON   := \033[0;36m
WHITE_ON    := \033[1;37m

BOLD := \033[1m
NORM := \033[21m

FAIL := $(RED_ON)$(BOLD)[FAIL]$(NORM)$(COLOR_FG)
OKAY := $(GREEN_ON)$(BOLD)[OKAY]$(NORM)$(COLOR_FG)
WARN := $(YELLOW_ON)$(BOLD)[WARN]$(NORM)$(COLOR_FG)
SKIP := $(BLUE_ON)$(BOLD)[SKIP]$(NORM)$(COLOR_FG)
INFO := $(BLUE_ON)$(BOLD)[INFO]$(NORM)$(COLOR_FG)
DEBG := $(PURPLE_ON)$(BOLD)[DEBG]$(NORM)$(COLOR_FG)

DOCKER_BIN     ?= docker

#--------------------------
# Drone-related configuration
#--------------------------
DRONE_ENV_FILE := .drone-env
DRONE_BIN      ?= drone
DRONE_EVENT    ?= push#:			Set the Drone CI event
DRONE_BRANCH   ?= $(shell git rev-parse --abbrev-ref HEAD)#: 		Set the Drone CI branch
DRONE_CMD       = $(DRONE_BIN) exec --trusted \
                --event $(DRONE_EVENT) \
                --branch $(DRONE_BRANCH) \
                --env-file $(DRONE_ENV_FILE) \
                --netrc-username $$login \
                --netrc-password $$password \
                --netrc-machine github.com

#--------------------------
# Builder images used by
# Drone pipeline
#-------------------------
BUILDER_GO          ?= artifactory.netskope.io/qe/citools/builder-go:2.6.4-alpine
BUILDER_PROTOBUF    ?= artifactory.netskope.io/qe/citools/builder-protobuf:2.6.4

# Help overrides; a project will extend these to personalize things.
HELP_EPILOGUE ?=\#\#====\\n\\nMake project is usually a good place to start.\\n

.DEFAULT_GOAL := help

.ONESHELL:
#--------------------------
# Help target, prints available targets
#--------------------------
.PHONY: help
help: ## This help output
	@printf "Available targets for $(PROJECT):\n\n"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/:.*##/:/'
	@printf "$(HELP_EPILOGUE)"

#--------------------------
# Display the optional variables.
#--------------------------
.PHONY: help-vars
help-vars::		## List of optional vars and their use
	@echo ""
	@echo "Optional variables you can configure as appropriate:"
	@echo ""
	@fgrep -h "#:" $(MAKEFILE_LIST) | sed -e 's/^# //g' | fgrep -v fgrep | fgrep -v my_help_text | sed -e 's/?.*#://' | sed -e 's/^export //g' | sort
	@echo ""

#---------------------
# Fetch the builder images from artifactory. When 'drone exec' runs a pipeline
# it uses these docker images. The first time drone fetches these images there is
# no feedback written to stdout during the fetch. It can appear stalled. This
# explicit target gives an opportunity to avoid a poor developer experience.
#---------------------
.PHONY: pull-builders
pull-builders:: env-check		## Pull the builder images from artifactory
	@# pull the builder images
	[ "$(VERBOSE)" = "yes" ] && {
		echo ">>-> Fetching builder images..."
	}
	docker pull -q $(BUILDER_GO)
	docker pull -q $(BUILDER_PROTOBUF)
	echo "$(OKAY) Builders fetched."

#---------------------------
# Set up the build environment, namely the env vars we want to inject
# into the pipeline, when running locally via 'drone exec'. Also
# creates the /tmp drone build cache location, which is shared for
# all/most drone steps
#---------------------------
.PHONY: env-setup
env-setup:: env-check		## Setup build environment
	@$(shell echo "UID=`id -u`" > $(DRONE_ENV_FILE))
	@$(shell echo "GID=`id -g`" >> $(DRONE_ENV_FILE))
	@$(shell echo "NS_ARTIFACTORY_HOST=artifactory.netskope.io" >> $(DRONE_ENV_FILE))
	@$(shell echo "DRONE_BUILD_NUMBER=dev" >> $(DRONE_ENV_FILE))
	@mkdir -p .drone-cache/go-build .drone-cache/go
	@chmod 755 scripts/getnrc

#----------------------
# Target to check if drone is installed. Many targets in the included make files
# will depend on drone.
#----------------------
env-check:: 	## Check that environment is sane for building
	@# environment checks
	command -v $(DOCKER_BIN) >/dev/null || {
		echo -e "$(FAIL) Could not find docker." 1>&2
		echo "       You need to install docker to build piratetreasure." 1>&2
		echo "       Goodbye!" 1>&2
		exit 1
	}
	command -v $(DRONE_BIN) >/dev/null || {
		echo -e "$(FAIL) Could not find your drone binary." 1>&2
		echo "       You need to install drone to build piratetreasure. Refer: https://docs.drone.io/cli/install/" 1>&2
		echo "       Goodbye!" 1>&2
		exit 1
	}
	cmd="chmod 755 scripts/getnrc"
	$$cmd || {
		echo -e "$(FAIL) Failure while changing scripts/getnrc permissions" 1>&2
		exit 1
	}
	exit 0
