# syntax=docker/dockerfile:1.4
########################################################################################################################
# Omnibus Image
########################################################################################################################

######## Build the frontend
FROM --platform=${BUILDPLATFORM} node AS frontendbuild
WORKDIR /go/src/github.com/analogj/scrutiny
COPY --link . /go/src/github.com/analogj/scrutiny

RUN make binary-frontend


######## Build the backend
FROM golang:1.25-bookworm as backendbuild

WORKDIR /go/src/github.com/analogj/scrutiny
COPY --link . /go/src/github.com/analogj/scrutiny
RUN apt-get update && DEBIAN_FRONTEND=noninteractive \
    apt-get install -y --no-install-recommends \
    file \
    && rm -rf /var/lib/apt/lists/*
RUN make binary-clean binary-all WEB_BINARY_NAME=scrutiny


######## Combine build artifacts in runtime image
FROM debian:bookworm-slim as runtime
ARG TARGETARCH
EXPOSE 8080
WORKDIR /opt/scrutiny
ENV PATH="/opt/scrutiny/bin:${PATH}"
ENV INFLUXD_CONFIG_PATH=/opt/scrutiny/influxdb
ENV S6VER="3.1.6.2"
ENV INFLUXVER="2.2.0"
ENV S6_SERVICES_READYTIME=1000
SHELL ["/usr/bin/sh", "-c"]

RUN apt-get update && DEBIAN_FRONTEND=noninteractive \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    cron \
    curl \
    smartmontools \
    tzdata \
    procps \
    xz-utils \
    && rm -rf /var/lib/apt/lists/* \
    && update-ca-certificates \
    && case ${TARGETARCH} in \
           "amd64")  S6_ARCH=x86_64  ;; \
           "arm64")  S6_ARCH=aarch64  ;; \
       esac \
    && curl https://github.com/just-containers/s6-overlay/releases/download/v${S6VER}/s6-overlay-noarch.tar.xz -L -s --output /tmp/s6-overlay-noarch.tar.xz \
    && tar -Jxpf /tmp/s6-overlay-noarch.tar.xz -C / \
    && rm -rf /tmp/s6-overlay-noarch.tar.xz \
    && curl https://github.com/just-containers/s6-overlay/releases/download/v${S6VER}/s6-overlay-${S6_ARCH}.tar.xz -L -s --output /tmp/s6-overlay-${S6_ARCH}.tar.xz \
    && tar -Jxpf /tmp/s6-overlay-${S6_ARCH}.tar.xz -C / \
    && rm -rf /tmp/s6-overlay-${S6_ARCH}.tar.xz
RUN curl -L https://dl.influxdata.com/influxdb/releases/influxdb2-${INFLUXVER}-${TARGETARCH}.deb --output /tmp/influxdb2-${INFLUXVER}-${TARGETARCH}.deb \
    && dpkg -i --force-all /tmp/influxdb2-${INFLUXVER}-${TARGETARCH}.deb \
    && rm -rf /tmp/influxdb2-${INFLUXVER}-${TARGETARCH}.deb

COPY /rootfs /

COPY --link --from=backendbuild --chmod=755 /go/src/github.com/analogj/scrutiny/scrutiny /opt/scrutiny/bin/
COPY --link --from=backendbuild --chmod=755 /go/src/github.com/analogj/scrutiny/scrutiny-collector-metrics /opt/scrutiny/bin/
COPY --link --from=frontendbuild --chmod=644 /go/src/github.com/analogj/scrutiny/dist /opt/scrutiny/web
RUN chmod 0644 /etc/cron.d/scrutiny && \
    rm -f /etc/cron.daily/* && \
    mkdir -p /opt/scrutiny/web && \
    mkdir -p /opt/scrutiny/config && \
    chmod -R ugo+rwx /opt/scrutiny/config && \
    chmod +x /etc/cont-init.d/* && \
    chmod +x /etc/services.d/*/run && \
    chmod +x /etc/services.d/*/finish

ENTRYPOINT ["/init"]
