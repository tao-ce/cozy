
build-image: check
	sudo buildah build \
		-t ${REGISTRY}/$(IMAGE_MEDIUM)/$(IMAGE_FULL_NAME):latest \
		--build-arg FEDORA_FLAVOR=$(BASE_IMAGE) \
		--build-arg FEDORA_VERSION=$(BASE_IMAGE_TAG) \
		--build-arg MEDIUM=$(IMAGE_MEDIUM) \
		--target appliance-$(IMAGE_MEDIUM) \
		--network=host \
 		.

$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).qcow2 \
$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).raw \
$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).vmdk:
	mkdir -p $(OUTPUT_DIR)
	sudo podman run \
		--rm \
		-it \
		--privileged \
		--security-opt label=type:unconfined_t \
		-v ./config/bootc.$(IMAGE_MEDIUM).toml:/config.toml:ro \
		-v $(OUTPUT_DIR):/output \
		-v  /var/lib/containers/storage:/var/lib/containers/storage \
		ghcr.io/osbuild/image-builder-cli:latest \
			build \
			--bootc-default-fs $(FS_TYPE) \
			--progress verbose \
			--output-name $(IMAGE_FULL_NAME) \
			--output-dir /output \
			--blueprint /config.toml \
			--bootc-ref $(REGISTRY)/$(IMAGE_MEDIUM)/$(IMAGE_FULL_NAME):latest \
			$(IMAGE_FORMAT)
	sudo chown $(shell id -u):$(shell id -g) $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).$(IMAGE_FORMAT)
