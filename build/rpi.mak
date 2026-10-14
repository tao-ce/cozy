ifeq ($(FS_TYPE),btrfs)
	ROOT_MOUNT_POINT="/:subvol=root"
else
	ROOT_MOUNT_POINT="/"
endif

build-rpi5-image: IMAGE_MEDIUM = rpi5
build-rpi5-image: build-image

build-rpi5: IMAGE_FORMAT = raw
build-rpi5: IMAGE_MEDIUM = rpi5
build-rpi5: $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).raw.xz

build-rpi4-image: IMAGE_MEDIUM = rpi4
build-rpi4-image: build-image

build-rpi4: IMAGE_FORMAT = raw
build-rpi4: IMAGE_MEDIUM = rpi4
build-rpi4: $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).raw.xz

$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).raw.xz: $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).raw
	sudo guestfish \
		-m /dev/sda3:$(ROOT_MOUNT_POINT) \
		-m /dev/sda2:/boot \
		-m /dev/sda1:/boot/efi \
		-a $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).$(IMAGE_FORMAT) \
		-f build/copy_efi_bootloader.fish

	ls -lh $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).$(IMAGE_FORMAT)
	xz -0 $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).$(IMAGE_FORMAT)
	ls -lh $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).$(IMAGE_FORMAT).xz
