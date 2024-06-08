# syntax = docker/dockerfile:1-experimental
FROM golang:1.20-alpine AS build

# Install dependencies
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash git openssh make build-base

WORKDIR /build

ADD attacker /build/attacker-service

#ADD https://api.github.com/repos/tsinghua-cel/attacker-service aversion.json
#RUN git clone -b {{ .Version }} --single-branch https://github.com/tsinghua-cel/attacker-service

RUN --mount=type=cache,target=/go/pkg/mod \
    cd /build/attacker-service && go mod download

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd /build/attacker-service && make

FROM alpine

RUN apk update && \
    apk upgrade && \
    apk add --no-cache build-base

WORKDIR /root

COPY  --from=build /build/attacker-service/build/bin/attacker /usr/bin/attacker
COPY ./entrypoint/attacker.sh /usr/local/bin/attacker.sh
RUN chmod u+x /usr/local/bin/attacker.sh

ENTRYPOINT [ "/usr/local/bin/attacker.sh" ]
