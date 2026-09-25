ifeq ($(FS_TYPE),btrfs)
	ROOT_MOUNT_POINT="/:subvol=root"
else
	ROOT_MOUNT_POINT="/"
endif

build-rpi4: IMAGE_FORMAT = raw
build-rpi4: $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).raw.xz

$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).raw.xz: $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).raw
	guestfish \
		-m /dev/sda3:${ROOT_MOUNT_POINT} \
		-m /dev/sda2:/boot \
		-m /dev/sda1:/boot/efi \
		-a $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).$(IMAGE_FORMAT) \
		-f build/copy_efi_bootloader.fish

	ls -lh $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).$(IMAGE_FORMAT)
	xz -0 $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).$(IMAGE_FORMAT)
	ls -lh $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).$(IMAGE_FORMAT).xz
	