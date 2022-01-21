#
# Copyright (c) 2022 Netskope, Inc. All rights reserved.
#
# Original Author: author <svulpius@netskope.com> (Jan 21, 2022)
#
# Description: its cursed
#
# Do NOT edit the following line
# GOPROJECT_VERSION=(devel)
#

# Store the location of this Makefile so we can use it as a path for finding
# other things relative to this location.
export MKFILE     := $(realpath $(lastword $(MAKEFILE_LIST)))
export MKFILE_DIR := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))

##==== Targets from piratetreasure/Makefile ====
include $(MKFILE_DIR)/build/mk/project-common.mk
include $(MKFILE_DIR)/build/mk/drone.mk
include $(MKFILE_DIR)/build/mk/piratetreasure.mk

# Turn up the volume on output if so desired (VERBOSE=yes)
export VERBOSE ?= no#:			Enable more build output (yes/no)

# The variables below define the application
APP_VERSION := v0.0.0
APP_SRC_DIR := ./cmd/$${app}

# Internal project-related defines
export PROJECT := piratetreasure
PROJECT_MODULE := $(shell grep ^module go.mod | awk '{print $$2}')

# Compile to a statically-linked binary or not. Static linking is very costly
# so it's nice to reserve it for final builds and not targets that would be
# part of a code/build/test/debug cycle.
STATIC ?= no#:				Compile to a statically-linked binary

# Other variables. Do not uncomment these; they are here for documentation
# purposes (and picked up by the 'help-vars' target). If you do not use these
# variables, remove them from this comment.
#
# PROJECT_VERSION ?#:		Set the version of the project
# GIT_SHA ?#:			Hash value of the git HEAD commit for this build
# BUILD_HOST ?#:			Host where the build occurs
# BUILT_BY ?#:			User who did the build
# BUILT_TIME ?#:			Time when build was done
