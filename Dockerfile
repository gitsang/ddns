FROM golang:1.22-bookworm AS builder

ARG VERSION=latest

ENV GO111MODULE=on

RUN apt update \
    && apt install -y --no-install-recommends git \
    && rm -rf /var/lib/apt/lists/* \
    && go install github.com/gitsang/ddns@${VERSION}

FROM debian:bookworm AS dist

LABEL maintainer="sang <sang.chen@outlook.com>"

RUN apt update \
    && apt install -y --no-install-recommends apt-transport-https ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*

RUN update-ca-certificates

COPY --from=builder /go/bin/ddns /usr/bin/ddns

ENTRYPOINT [ "/usr/bin/ddns" ]
