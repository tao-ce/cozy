ARG FEDORA_VERSION=44
ARG FEDORA_FLAVOR=quay.io/fedora-ostree-desktops/kinoite
ARG FLAVOR=vm

ARG TAO_DOMAIN="tao-community-edition.local"

ARG OS_VARIANT="TAO Community Edition - Cozy"
ARG OS_VARIANT_ID="com.taotesting.cozy"


FROM ${FEDORA_FLAVOR}:${FEDORA_VERSION} AS appliance

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG FEDORA_VERSION
ARG FEDORA_FLAVOR
ARG OS_VARIANT
ARG OS_VARIANT_ID
ARG TAO_DOMAIN
ARG TIMEZONE="UTC"
ARG FLAVOR
ARG COCKPIT_TAO_CE_VERSION='0.0.0'

LABEL org.opencontainers.image.authors="opensource-support@taotesting.com"

# RUN dnf -y swap fedora-release generic-release --allowerasing \
#     && dnf -y remove kdeconnectd  kde-connect \
#     && echo 'VARIANT="${OS_VARIANT}"' >>/usr/lib/os-release \
#     && echo 'VARIANT_ID="${OS_VARIANT_ID}"' >>/usr/lib/os-release

# COPY ./config/images.common.lst /run/context/images.common.lst
# COPY ./config/images.${FLAVOR}.lst /run/context/images.${FLAVOR}.lst
# RUN \
#     cat /run/context/images.*.lst \
#         | grep -v '^#' \
#         | grep -v '^[ ]*$' \
#         | xargs -n1 \
#             podman pull

# COPY ./config/packages.common.lst /run/context/packages.common.lst
# COPY ./config/packages.${FLAVOR}.lst /run/context/packages.${FLAVOR}.lst

# ADD https://copr.fedorainfracloud.org/coprs/pbrobinson/a64-kernel/repo/fedora-${FEDORA_VERSION}/pbrobinson-a64-kernel-fedora-${FEDORA_VERSION}.repo /etc/yum.repos.d/pbrobinson-a64-kernel-fedora-${FEDORA_VERSION}.repo
ADD \
    https://copr.fedorainfracloud.org/coprs/dwrobel/kernel-rpi/repo/fedora-${FEDORA_VERSION}/dwrobel-kernel-rpi-fedora-${FEDORA_VERSION}.repo \
    /etc/yum.repos.d/dwrobel-kernel-rpi-fedora-${FEDORA_VERSION}.repo
ADD \
    https://copr.fedorainfracloud.org/coprs/dwrobel/bcm283x-firmware-rpi/repo/fedora-${FEDORA_VERSION}/dwrobel-bcm283x-firmware-rpi-fedora-${FEDORA_VERSION}.repo \
    /etc/yum.repos.d/dwrobel-bcm283x-firmware-rpi-fedora-${FEDORA_VERSION}.repo
ADD \
    https://copr.fedorainfracloud.org/coprs/dwrobel/bcm434xx-firmware-rpi/repo/fedora-${FEDORA_VERSION}/dwrobel-bcm434xx-firmware-rpi-fedora-${FEDORA_VERSION}.repo \
    /etc/yum.repos.d/dwrobel-bcm434xx-firmware-rpi-fedora-${FEDORA_VERSION}.repo

RUN ls -lrth /usr/lib/modules/* /boot /usr/lib/ostree-boot /boot/efi /usr/lib/ostree-boot/efi || true

RUN \
    --mount=type=cache,id=dnf-cache,target=/var/cache/dnf \
    --mount=type=cache,id=libdnf-cache,target=/var/cache/libdnf5 \
    rpm-ostree override replace\
        --freeze \
        --experimental \
        --from repo='copr:copr.fedorainfracloud.org:dwrobel:kernel-rpi' \
        --remove kernel-modules-core \
        kernel kernel-{core,modules} 

RUN ls -lrth /usr/lib/modules/* /boot /usr/lib/ostree-boot /boot/efi /usr/lib/ostree-boot/efi || true
    
# RUN \
#     --mount=type=cache,id=dnf-cache,target=/var/cache/dnf \
#     --mount=type=cache,id=libdnf-cache,target=/var/cache/libdnf5 \
#     cat /run/context/packages.*.lst \
#         | grep -v '^#' \
#         | grep -v '^[ ]*$' \
#         | xargs \
#             dnf \
#                 --setopt=install_weak_deps=false \
#                 install -y
# # 
# COPY root/ /
# ADD https://github.com/tao-ce/cockpit-tao-ce/releases/download/${COCKPIT_TAO_CE_VERSION}/cockpit-tao-ce-${COCKPIT_TAO_CE_VERSION}.tar.xz /run/tmp/cockpit-tao-ce.tar.xz
# RUN mkdir -p /usr/share/cockpit/tao-ce/ \
    # && tar -xf \
        # /run/tmp/cockpit-tao-ce.tar.xz \
        # -C /usr/share/cockpit/tao-ce/ \
        # --strip-components=2 \
        # cockpit-tao-ce/dist
# 
RUN \
    ln -s /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    && systemctl enable \
        sshd.service \
        avahi-daemon.service \
    && dnf clean all
