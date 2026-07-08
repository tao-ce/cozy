ARG FEDORA_VERSION=44
ARG FEDORA_FLAVOR=quay.io/fedora-ostree-desktops/kinoite
ARG MEDIUM=vm

ARG TAO_DOMAIN="tao-community-edition.local"

ARG OS_VARIANT="TAO Community Edition - Cozy"
ARG OS_VARIANT_ID="com.taotesting.cozy"


FROM golang:1-alpine AS appliance-setup-builder
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY ./packages/appliance-setup/ .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-s -w" -o appliance-setup main.go

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
ARG MEDIUM
ARG COCKPIT_TAO_CE_VERSION='0.0.0'

LABEL org.opencontainers.image.authors="opensource-support@taotesting.com"

RUN dnf -y swap fedora-release generic-release --allowerasing \
    && dnf -y remove kdeconnectd  kde-connect \
    && echo 'VARIANT="${OS_VARIANT}"' >>/usr/lib/os-release \
    && echo 'VARIANT_ID="${OS_VARIANT_ID}"' >>/usr/lib/os-release

COPY ./config/images.common.lst /run/context/images.common.lst
COPY ./config/images.${MEDIUM}.lst /run/context/images.${MEDIUM}.lst
RUN \
    cat /run/context/images.*.lst \
        | grep -v '^#' \
        | grep -v '^[ ]*$' \
        | xargs -n1 \
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

COPY root/ /
ADD https://github.com/tao-ce/cockpit-tao-ce/releases/download/${COCKPIT_TAO_CE_VERSION}/cockpit-tao-ce-${COCKPIT_TAO_CE_VERSION}.tar.xz /run/tmp/cockpit-tao-ce.tar.xz
RUN mkdir -p /usr/share/cockpit/tao-ce/ \
    && tar -xf \
        /run/tmp/cockpit-tao-ce.tar.xz \
        -C /usr/share/cockpit/tao-ce/ \
        --strip-components=2 \
        cockpit-tao-ce/dist

COPY --from=appliance-setup-builder \
    /app/appliance-setup \
    /usr/libexec/tao-ce/cozy/appliance-setup

RUN \
    mkdir -p /etc/tao-ce/cozy \
    && echo ${MEDIUM} >/etc/tao-ce/cozy/medium \
    && ln -s /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    && ln -s /usr/lib/systemd/system/kmsconvt@.service /etc/systemd/system/autovt@.service \
    && systemctl disable getty@tty1.service \
    && systemctl enable \
        sshd.service \
        appliance-setup@tty4.service \
        kmscon.service \
    && dnf clean all
