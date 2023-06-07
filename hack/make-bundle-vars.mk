# Image URL to use all building/pushing image targets
IMAGE_REGISTRY ?= quay.io
REGISTRY_NAMESPACE ?= jarrpa
DEFAULT_IMAGE_TAG ?= latest

REPORTER_IMAGE_NAME ?= managed-fusion-metering-reporter
REPORTER_IMAGE_TAG ?= latest

# *_IMG variables define the final image names
REPORTER_IMG ?= $(IMAGE_REGISTRY)/$(REGISTRY_NAMESPACE)/$(REPORTER_IMAGE_NAME):$(REPORTER_IMAGE_TAG)
