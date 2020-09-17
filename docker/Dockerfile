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
COPY ../webapp/frontend /scrutiny/src

RUN npm install -g @angular/cli@9.1.4 && \
    mkdir -p /scrutiny/dist && \
    npm install && \
    ng build --output-path=/scrutiny/dist --deploy-url="/web/" --base-href="/web/" --prod


########
FROM ubuntu:bionic as runtime
EXPOSE 8080
WORKDIR /scrutiny
ENV PATH="/scrutiny/bin:${PATH}"

ADD https://github.com/dshearer/jobber/releases/download/v1.4.4/jobber_1.4.4-1_amd64.deb /tmp/
RUN apt install /tmp/jobber_1.4.4-1_amd64.deb

RUN apt-get update && apt-get install -y smartmontools=7.0-0ubuntu1~ubuntu18.04.1

ADD https://github.com/just-containers/s6-overlay/releases/download/v1.21.8.0/s6-overlay-amd64.tar.gz /tmp/
RUN tar xzf /tmp/s6-overlay-amd64.tar.gz -C /
COPY /rootfs /


COPY --from=backendbuild /go/src/github.com/analogj/scrutiny/scrutiny /scrutiny/bin/
COPY --from=backendbuild /go/src/github.com/analogj/scrutiny/scrutiny-collector-selftest /scrutiny/bin/
COPY --from=backendbuild /go/src/github.com/analogj/scrutiny/scrutiny-collector-metrics /scrutiny/bin/
COPY --from=frontendbuild /scrutiny/dist /scrutiny/web
RUN chmod +x /scrutiny/bin/scrutiny && \
    chmod +x /scrutiny/bin/scrutiny-collector-selftest && \
    chmod +x /scrutiny/bin/scrutiny-collector-metrics && \
    mkdir -p /scrutiny/web && \
    mkdir -p /scrutiny/config && \
    mkdir -p /scrutiny/jobber

CMD ["/init"]
