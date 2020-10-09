########
FROM golang:1.14.4-buster as backendbuild

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

WORKDIR /scrutiny/src
COPY webapp/frontend /scrutiny/src

RUN npm install -g @angular/cli@9.1.4 && \
    mkdir -p /scrutiny/dist && \
    npm install && \
    ng build --output-path=/scrutiny/dist --prod


########
FROM ubuntu:bionic as runtime
EXPOSE 8080
WORKDIR /scrutiny
ENV PATH="/scrutiny/bin:${PATH}"

RUN apt-get update && apt-get install -y cron smartmontools=7.0-0ubuntu1~ubuntu18.04.1 ca-certificates curl && update-ca-certificates

ADD https://github.com/just-containers/s6-overlay/releases/download/v1.21.8.0/s6-overlay-amd64.tar.gz /tmp/
RUN tar xzf /tmp/s6-overlay-amd64.tar.gz -C /
COPY /rootfs /

COPY /rootfs/etc/cron.d/scrutiny /etc/cron.d/scrutiny
COPY --from=backendbuild /go/src/github.com/analogj/scrutiny/scrutiny /scrutiny/bin/
COPY --from=backendbuild /go/src/github.com/analogj/scrutiny/scrutiny-collector-selftest /scrutiny/bin/
COPY --from=backendbuild /go/src/github.com/analogj/scrutiny/scrutiny-collector-metrics /scrutiny/bin/
COPY --from=frontendbuild /scrutiny/dist /scrutiny/web
RUN chmod +x /scrutiny/bin/scrutiny && \
    chmod +x /scrutiny/bin/scrutiny-collector-selftest && \
    chmod +x /scrutiny/bin/scrutiny-collector-metrics && \
    chmod 0644 /etc/cron.d/scrutiny && \
    rm -f /etc/cron.daily/* && \
    mkdir -p /scrutiny/web && \
    mkdir -p /scrutiny/config && \
    chmod -R ugo+rwx /scrutiny/config


CMD ["/init"]
