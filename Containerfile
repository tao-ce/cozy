ARG FEDORA_VERSION=44
ARG FEDORA_FLAVOR=quay.io/fedora-ostree-desktops/kinoite
ARG FLAVOR=vm

ARG TAO_DOMAIN="tao-community-edition.local"

ARG OS_VARIANT="TAO Community Edition - Cozy"
ARG OS_VARIANT_ID="com.taotesting.cozy"

FROM ${FEDORA_FLAVOR}:${FEDORA_VERSION} AS base

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG FEDORA_VERSION
ARG FEDORA_FLAVOR
ARG OS_VARIANT
ARG OS_VARIANT_ID

RUN dnf -y swap fedora-release generic-release --allowerasing \
    && dnf -y remove kdeconnectd  kde-connect \
    && echo 'VARIANT="${OS_VARIANT}"' >>/usr/lib/os-release \
    && echo 'VARIANT_ID="${OS_VARIANT_ID}"' >>/usr/lib/os-release

FROM node:24 AS build-cockpit

ENV NODE_ENV=production
WORKDIR /app
COPY src/cockpit-tao-ce/ /app/
RUN \
    --mount=type=cache,id=npm-cache,target=/root/.npm,sharing=locked \
    make

FROM base AS appliance

LABEL org.opencontainers.image.authors="opensource-support@taotesting.com"

ARG TAO_DOMAIN
ARG TIMEZONE="UTC"
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG FLAVOR

COPY ./config/images.lst /run/context/images.lst
COPY ./config/images.${FLAVOR}.lst /run/context/images.${FLAVOR}.lst
RUN \
    cat /run/context/images.*.lst \
        | grep -v '^#' \
        | grep -v '^[ ]*$' \
        | xargs -n1 \
            podman pull

COPY ./config/packages.common.lst /run/context/packages.common.lst
COPY ./config/packages.${FLAVOR}.lst /run/context/packages.${FLAVOR}.lst

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
COPY --from=build-cockpit /app/dist/ /usr/share/cockpit/tao-ce/

RUN \
    ln -s /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    && systemctl enable \
        sshd.service \
        avahi-daemon.service \
    && dnf clean all
