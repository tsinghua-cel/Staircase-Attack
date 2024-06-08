# syntax = docker/dockerfile:1-experimental
FROM golang:1.20-alpine AS build

# Install dependencies
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash git openssh make build-base

RUN go env -w CGO_ENABLED="1"

WORKDIR /build

ADD prysm_modified /build/prysm

#ADD https://api.github.com/repos/tsinghua-cel/prysm pversion.json
#RUN git clone -b {{ .Version }}  --single-branch https://github.com/tsinghua-cel/prysm.git

RUN --mount=type=cache,target=/go/pkg/mod \
    cd /build/prysm && go mod download

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd /build/prysm && go build -o /beacon-chain ./cmd/beacon-chain

FROM alpine

RUN apk update && \
    apk upgrade && \
    apk add --no-cache build-base

WORKDIR /root

COPY  --from=build /beacon-chain /usr/bin/beacon-chain
COPY ./entrypoint/beacon-node.sh /usr/local/bin/beacon-node.sh
RUN chmod u+x /usr/local/bin/beacon-node.sh

ENTRYPOINT [ "/usr/local/bin/beacon-node.sh" ]
