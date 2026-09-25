REGISTRY := ghcr.io/tao-ce/cozy
BASE_IMAGE := quay.io/fedora/fedora-bootc
BASE_IMAGE_TAG := 44
FS_TYPE := btrfs

IMAGE_TYPE := vm
IMAGE_NAME := server

IMAGE_FULL_NAME := tao-ce-cozy-$(IMAGE_TYPE)-$(IMAGE_NAME)
RUNNER_DISTRO := $(shell grep -oP '^ID=\K[^"]+' /etc/os-release)

OUTPUT_DIR := ./output

requirements: requirements-$(RUNNER_DISTRO)

requirements-ubuntu:
	sudo apt update
	sudo apt install -y buildah jq jsonnet qemu-utils libguestfs-tools

requirements-fedora:
	sudo dnf install -y buildah jq jsonnet qemu-img libguestfs-tools

include build/bootc.mak
include build/vm.mak
include build/rpi4.mak

check: check-env

check-env:
ifndef REGISTRY
	$(error REGISTRY is not set)
endif
ifndef BASE_IMAGE
	$(error BASE_IMAGE is not set)
endif
ifndef BASE_IMAGE_TAG
	$(error BASE_IMAGE_TAG is not set)
endif
ifndef FS_TYPE
	$(error FS_TYPE is not set)
endif
ifndef OUTPUT_DIR
	$(error OUTPUT_DIR is not set)
endif
ifndef IMAGE_NAME
	$(error IMAGE_NAME is not set)
endif
ifndef IMAGE_TYPE
	$(error IMAGE_TYPE is not set)
endif