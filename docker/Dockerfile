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
FROM golang:1.20-bookworm as backendbuild

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
ENV S6VER="1.21.8.0"
ENV INFLUXVER="2.2.0"
SHELL ["/usr/bin/sh", "-c"]

RUN apt-get update && DEBIAN_FRONTEND=noninteractive \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    cron \
    curl \
    smartmontools \
    tzdata \
    && rm -rf /var/lib/apt/lists/* \
    && update-ca-certificates \
    && case ${TARGETARCH} in \
           "amd64")  S6_ARCH=amd64  ;; \
           "arm64")  S6_ARCH=aarch64  ;; \
       esac \
    && curl https://github.com/just-containers/s6-overlay/releases/download/v${S6VER}/s6-overlay-${S6_ARCH}.tar.gz -L -s --output /tmp/s6-overlay-${S6_ARCH}.tar.gz \
    && tar xzf /tmp/s6-overlay-${S6_ARCH}.tar.gz -C / \
    && rm -rf /tmp/s6-overlay-${S6_ARCH}.tar.gz \
    curl -L https://dl.influxdata.com/influxdb/releases/influxdb2-${INFLUXVER}-${TARGETARCH}.deb --output /tmp/influxdb2-${INFLUXVER}-${TARGETARCH}.deb \
    && ln -s /usr/bin/false /bin/false \
    && ln -s /usr/bin/bash /bin/bash \
    && dpkg -i --force-all /tmp/influxdb2-${INFLUXVER}-${TARGETARCH}.deb \
    && rm -f /bin/bash \
    && rm -rf /tmp/influxdb2-${INFLUXVER}-${TARGETARCH}.deb

COPY /rootfs /

COPY --link --from=backendbuild --chmod=755 /go/src/github.com/analogj/scrutiny/scrutiny /opt/scrutiny/bin/
COPY --link --from=backendbuild --chmod=755 /go/src/github.com/analogj/scrutiny/scrutiny-collector-metrics /opt/scrutiny/bin/
COPY --link --from=frontendbuild --chmod=644 /go/src/github.com/analogj/scrutiny/dist /opt/scrutiny/web
RUN chmod 0644 /etc/cron.d/scrutiny && \
    rm -f /etc/cron.daily/* && \
    mkdir -p /opt/scrutiny/web && \
    mkdir -p /opt/scrutiny/config && \
    chmod -R ugo+rwx /opt/scrutiny/config

CMD ["/init"]
