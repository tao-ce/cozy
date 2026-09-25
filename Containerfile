ARG FEDORA_VERSION=45
ARG FEDORA_FLAVOR=quay.io/fedora-ostree-desktops/kinoite
ARG MEDIUM=vm

ARG OS_VARIANT="TAO Community Edition - Cozy"
ARG OS_VARIANT_ID="com.taotesting.cozy"

<<<<<<< HEAD
=======

>>>>>>> 7370933 (feat: taosetup dashboard WIP)
FROM golang:1-alpine AS taosetup-builder
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY ./packages/taosetup/ .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-s -w" -o taosetup main.go

FROM ${FEDORA_FLAVOR}:${FEDORA_VERSION} AS base

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG FEDORA_VERSION
ARG FEDORA_FLAVOR
ARG OS_VARIANT
ARG OS_VARIANT_ID
ARG TIMEZONE="UTC"
ARG MEDIUM
ARG COCKPIT_TAO_CE_VERSION='0.0.0'

LABEL org.opencontainers.image.authors="opensource-support@taotesting.com"

RUN \
    --mount=type=cache,id=dnf-cache,target=/var/cache/dnf \
    --mount=type=cache,id=libdnf-cache,target=/var/cache/libdnf5 \
    dnf -y swap fedora-release generic-release --allowerasing \
        && dnf -y remove kdeconnectd  kde-connect \
        && echo 'VARIANT="${OS_VARIANT}"' >>/usr/lib/os-release \
        && echo 'VARIANT_ID="${OS_VARIANT_ID}"' >>/usr/lib/os-release

COPY ./config/images.common.lst /run/context/images.common.lst
COPY ./config/images.${MEDIUM}.lst /run/context/images.${MEDIUM}.lst
RUN \
    cat /run/context/images.*.lst \
        | grep -v '^#' \
        | grep -v '^[ ]*$' \
        | xargs -n2 \
            podman pull

COPY ./config/packages.common.lst /run/context/packages.common.lst
COPY ./config/packages.${MEDIUM}.lst /run/context/packages.${MEDIUM}.lst

RUN \
    --mount=type=cache,id=dnf-cache,target=/var/cache/dnf \
    --mount=type=cache,id=libdnf-cache,target=/var/cache/libdnf5 \
    cat /run/context/packages.*.lst \
        | grep -v '^#' \
        | grep -v '^[ ]*$' \
        | xargs \
            dnf \
                --setopt=install_weak_deps=false \
                install -y

ADD https://github.com/tao-ce/cockpit-tao-ce/releases/download/${COCKPIT_TAO_CE_VERSION}/cockpit-tao-ce-${COCKPIT_TAO_CE_VERSION}.tar.xz /run/tmp/cockpit-tao-ce.tar.xz
RUN mkdir -p /usr/share/cockpit/tao-ce/ \
    && tar -xf \
        /run/tmp/cockpit-tao-ce.tar.xz \
        -C /usr/share/cockpit/tao-ce/ \
        --strip-components=2 \
        cockpit-tao-ce/dist

COPY --from=taosetup-builder \
    /app/taosetup \
    /usr/libexec/tao-ce/cozy/taosetup

COPY root/ /

RUN \
    mkdir -p /etc/tao-ce/cozy \
    && echo ${MEDIUM} >/etc/tao-ce/cozy/medium \
    && ln -s /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    && rm -f /etc/systemd/system/autovt@.service \
    && ln -s /usr/lib/systemd/system/kmsconvt@.service /etc/systemd/system/autovt@.service \
    && systemctl disable getty@tty1.service \
    && systemctl enable \
        sshd.service \
        taosetup-install@tty4.service \
        taosetup-dashboard@tty4.service \
        kmscon.service \
    && dnf clean all

FROM base AS appliance-vm
FROM base AS appliance-rpi4
FROM base AS appliance-rpi5

ARG FEDORA_VERSION

ADD \
    https://copr.fedorainfracloud.org/coprs/dwrobel/kernel-rpi/repo/fedora-${FEDORA_VERSION}/dwrobel-kernel-rpi-fedora-${FEDORA_VERSION}.repo \
    /etc/yum.repos.d/dwrobel-kernel-rpi-fedora-${FEDORA_VERSION}.repo
ADD \
    https://copr.fedorainfracloud.org/coprs/dwrobel/bcm283x-firmware-rpi/repo/fedora-${FEDORA_VERSION}/dwrobel-bcm283x-firmware-rpi-fedora-${FEDORA_VERSION}.repo \
    /etc/yum.repos.d/dwrobel-bcm283x-firmware-rpi-fedora-${FEDORA_VERSION}.repo
ADD \
    https://copr.fedorainfracloud.org/coprs/dwrobel/bcm434xx-firmware-rpi/repo/fedora-${FEDORA_VERSION}/dwrobel-bcm434xx-firmware-rpi-fedora-${FEDORA_VERSION}.repo \
    /etc/yum.repos.d/dwrobel-bcm434xx-firmware-rpi-fedora-${FEDORA_VERSION}.repo

RUN \
    --mount=type=cache,id=dnf-cache,target=/var/cache/dnf \
    --mount=type=cache,id=libdnf-cache,target=/var/cache/libdnf5 \
    rpm-ostree override replace\
        --freeze \
        --experimental \
        --from repo='copr:copr.fedorainfracloud.org:dwrobel:kernel-rpi' \
        --remove kernel \
        --remove kernel-core \
        --remove kernel-modules \
        --remove kernel-modules-core \
        kernel-rpi5 \
        && cd /usr/lib/modules/*/ && gzip -c vmlinux >vmlinuz  \
        && dracut --force initramfs.img "$(expr $(pwd) : '.*[/]\([^/]*\)$')" \
        && dnf clean all
