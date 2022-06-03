########
FROM golang:1.17.10-buster as backendbuild

WORKDIR /go/src/github.com/analogj/scrutiny

COPY . /go/src/github.com/analogj/scrutiny

RUN go mod vendor && \
    go build -ldflags '-w -extldflags "-static"' -o scrutiny webapp/backend/cmd/scrutiny/scrutiny.go && \
    go build -ldflags '-w -extldflags "-static"' -o scrutiny-collector-selftest collector/cmd/collector-selftest/collector-selftest.go && \
    go build -ldflags '-w -extldflags "-static"' -o scrutiny-collector-metrics collector/cmd/collector-metrics/collector-metrics.go

########
FROM node:lts-slim as frontendbuild

#reduce logging, disable angular-cli analytics for ci environment
ENV NPM_CONFIG_LOGLEVEL=warn NG_CLI_ANALYTICS=false

WORKDIR /opt/scrutiny/src
COPY webapp/frontend /opt/scrutiny/src

RUN npm install -g @angular/cli@9.1.4 && \
    mkdir -p /scrutiny/dist && \
    npm install && \
    npm run build:prod -- --output-path=/opt/scrutiny/dist


########
FROM ubuntu:latest as runtime
ARG TARGETARCH
EXPOSE 8080
WORKDIR /opt/scrutiny
ENV PATH="/opt/scrutiny/bin:${PATH}"
ENV INFLUXD_CONFIG_PATH=/opt/scrutiny/influxdb

RUN apt-get update && apt-get install -y cron smartmontools ca-certificates curl tzdata \
    && update-ca-certificates \
    &&  case ${TARGETARCH} in \
            "amd64")  S6_ARCH=amd64  ;; \
            "arm64")  S6_ARCH=aarch64  ;; \
        esac \
    && curl https://github.com/just-containers/s6-overlay/releases/download/v1.21.8.0/s6-overlay-${S6_ARCH}.tar.gz -L -s --output /tmp/s6-overlay-${S6_ARCH}.tar.gz \
    && tar xzf /tmp/s6-overlay-${S6_ARCH}.tar.gz -C / --exclude="./bin" \
    && tar xzf /tmp/s6-overlay-${S6_ARCH}.tar.gz -C /usr ./bin \
    && rm -rf /tmp/s6-overlay-${S6_ARCH}.tar.gz \
    && curl -L https://dl.influxdata.com/influxdb/releases/influxdb2-2.2.0-${TARGETARCH}.deb --output /tmp/influxdb2-2.2.0-${TARGETARCH}.deb \
    && dpkg -i --force-all /tmp/influxdb2-2.2.0-${TARGETARCH}.deb

COPY /rootfs /

COPY /rootfs/etc/cron.d/scrutiny /etc/cron.d/scrutiny
COPY --from=backendbuild /go/src/github.com/analogj/scrutiny/scrutiny /opt/scrutiny/bin/
COPY --from=backendbuild /go/src/github.com/analogj/scrutiny/scrutiny-collector-selftest /opt/scrutiny/bin/
COPY --from=backendbuild /go/src/github.com/analogj/scrutiny/scrutiny-collector-metrics /opt/scrutiny/bin/
COPY --from=frontendbuild /opt/scrutiny/dist /opt/scrutiny/web
RUN chmod +x /opt/scrutiny/bin/scrutiny && \
    chmod +x /opt/scrutiny/bin/scrutiny-collector-selftest && \
    chmod +x /opt/scrutiny/bin/scrutiny-collector-metrics && \
    chmod 0644 /etc/cron.d/scrutiny && \
    rm -f /etc/cron.daily/* && \
    mkdir -p /opt/scrutiny/web && \
    mkdir -p /opt/scrutiny/config && \
    chmod -R ugo+rwx /opt/scrutiny/config


CMD ["/init"]
