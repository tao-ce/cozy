
build-image: check
	sudo buildah build \
		-t ${REGISTRY}/$(IMAGE_TYPE)/$(IMAGE_FULL_NAME):latest \
		--build-arg FEDORA_FLAVOR=$(BASE_IMAGE) \
		--build-arg FEDORA_VERSION=$(BASE_IMAGE_TAG) \
		--build-arg MEDIUM=$(IMAGE_TYPE) \
		--target appliance \
		--network=host \
# 		.
# $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).$(IMAGE_FORMAT):
# 	mkdir -p $(OUTPUT_DIR)/$(IMAGE_NAME)
# 	sudo podman run \
# 		--rm \
# 		-it \
# 		--privileged \
# 		--security-opt label=type:unconfined_t \
# 		-v ./config/bootc.$(IMAGE_TYPE).toml:/config.toml:ro \
# 		-v $(OUTPUT_DIR)/$(IMAGE_NAME):/output \
# 		-v  /var/lib/containers/storage:/var/lib/containers/storage \
# 		quay.io/centos-bootc/bootc-image-builder:latest \
# 			--type $(IMAGE_FORMAT) \
# 			--rootfs $(FS_TYPE) \
# 			--chown $(shell id -u):$(shell id -g) \
# 			--log-level debug \
# 			-vv \
# 			--use-librepo \
# 			$(REGISTRY)/$(IMAGE_TYPE)/$(IMAGE_FULL_NAME):latest
# 	mv $(OUTPUT_DIR)/$(IMAGE_NAME)/$(IMAGE_FORMAT)/disk.$(IMAGE_FORMAT) $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).$(IMAGE_FORMAT)


$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).qcow2 \
$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).raw \
$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).vmdk:
	sudo podman run \
		--rm \
		-it \
		--privileged \
		--security-opt label=type:unconfined_t \
		-v ./config/bootc.$(IMAGE_TYPE).toml:/config.toml:ro \
		-v $(OUTPUT_DIR):/output \
		-v  /var/lib/containers/storage:/var/lib/containers/storage \
		ghcr.io/osbuild/image-builder-cli:latest \
			build \
			--bootc-default-fs $(FS_TYPE) \
			--progress verbose \
			--output-name $(IMAGE_FULL_NAME) \
			--output-dir /output \
			--blueprint /config.toml \
			--bootc-ref $(REGISTRY)/$(IMAGE_TYPE)/$(IMAGE_FULL_NAME):latest \
			$(IMAGE_FORMAT)