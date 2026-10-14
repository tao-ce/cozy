REGISTRY = ghcr.io/tao-ce/cozy
BASE_IMAGE = quay.io/fedora/fedora-bootc
BASE_IMAGE_TAG = 45
FS_TYPE = btrfs

IMAGE_MEDIUM = vm
IMAGE_NAME = server

IMAGE_NAME_PREFIX = tao-ce-cozy
IMAGE_FULL_NAME = $(IMAGE_NAME_PREFIX)-$(IMAGE_MEDIUM)-$(IMAGE_NAME)
RUNNER_DISTRO := $(shell grep -oP '^ID=\K[^"]+' /etc/os-release)
RUNNER_VARIANT := $(shell grep -oP '^VARIANT_ID=\K[^"]+' /etc/os-release)

OUTPUT_DIR := ./output

ALLOWED_MEDIA := vm rpi4 rpi5
ifeq ($(filter $(IMAGE_MEDIUM), $(ALLOWED_MEDIA)),)
	$(error IMAGE_MEDIUM is not set to a valid medium. Allowed values are: $(ALLOWED_MEDIA))
endif

requirements: requirements-$(RUNNER_DISTRO)

requirements-ubuntu: requirements-debian

requirements-debian:
	sudo apt update
	sudo apt install -y buildah jq jsonnet qemu-utils libguestfs-tools

requirements-fedora:
ifeq ($(RUNNER_VARIANT), coreos)
	sudo rpm-ostree install -Ay --idempotent \
		buildah jq jsonnet qemu-img libguestfs-tools libvirt-daemon libvirt-daemon-driver-qemu
	sudo systemctl enable --now  virtqemud
else
	sudo dnf install -y buildah jq jsonnet qemu-img libguestfs-tools
endif

include build/bootc.mak
include build/vm.mak
include build/rpi.mak

check: check-env #check-group

check-group:
	groups | grep -q libvirt || (echo "You need to be in the libvirt group to build images"; exit 1)

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
ifndef IMAGE_MEDIUM
	$(error IMAGE_MEDIUM is not set)
endif