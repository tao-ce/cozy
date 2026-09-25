TEST_VM_IP := 127.0.47.101

build-ova: IMAGE_FORMAT = vmdk
build-ova: $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).ova

local-run: IMAGE_FORMAT = qcow2
local-run: $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).qcow2
	virt-install \
		--connect qemu:///system \
	    -n tao-ce.cozy \
		--cpu host \
		--memory 8192 \
		--vcpus 4 \
		--osinfo fedora-unknown \
		--import \
		--disk path=$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).qcow2,format=qcow2 \
		--network passt,portForward0=$(TEST_VM_IP):443:443,portForward1=$(TEST_VM_IP):8443:8443,portForward2=$(TEST_VM_IP):8022:22



$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).ovf: $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).vmdk
	jq -nc  \
		--arg disk_file $(IMAGE_FULL_NAME).vmdk \
		--arg disk_capacity $(shell qemu-img info --output json $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).vmdk | jq -r '.["virtual-size"]') \
		--arg disk_uuid $(shell uuidgen) \
		--arg machine_uuid $(shell uuidgen) \
		'{disk: {file: $$disk_file, capacity: $$disk_capacity, uuid: $$disk_uuid }, machine: { uuid: $$machine_uuid } }' \
		| jsonnet -S \
			--tla-code-file args=/dev/stdin \
			config/vm.ovf.jsonnet \
				> $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).ovf

$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).mf: $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).ovf $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).vmdk
	sha256sum --tag \
		$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).ovf \
		$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).vmdk \
		> $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).mf

$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).ova: $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).ovf $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).vmdk $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).mf
	tar -cvf $(OUTPUT_DIR)/$(IMAGE_FULL_NAME).ova \
		$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).ovf \
		$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).vmdk \
		$(OUTPUT_DIR)/$(IMAGE_FULL_NAME).mf

