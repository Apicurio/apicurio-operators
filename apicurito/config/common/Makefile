# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

VERSION ?= 1.0.1
PREVIOUS_VERSION := 0.1.0
VERSION ?= $(DEFAULT_VERSION)
DEFAULT_IMAGE ?=  quay.io/apicurio/apicurito-operator
IMAGE ?= $(DEFAULT_IMAGE)
DEFAULT_TAG := $(VERSION)
TAG ?= $(DEFAULT_TAG)
DEFAULT_CHANNEL := latest
CHANNEL ?= $(DEFAULT_CHANNEL)
DEFAULT_CHANNELS := $(DEFAULT_CHANNEL),candidate,stable,$(shell v=$(VERSION) && echo "$${v%\.[0-9]}.x")
CHANNELS ?= $(DEFAULT_CHANNELS)

CONTROLLER_GEN_VERSION := v0.4.1
OPERATOR_SDK_VERSION := v1.14.0
KUSTOMIZE_VERSION := v4.1.2

#
# Make Options
#
MK_OPTIONS := -s
