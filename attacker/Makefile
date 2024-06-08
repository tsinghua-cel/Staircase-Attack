.PHONY: default attacker reward all clean docker docs

GOBIN = $(shell pwd)/build/bin
TAG ?= latest
GOFILES_NOVENDOR := $(shell go list -f "{{.Dir}}" ./...)

VERSION := $(shell git describe --tags)
COMMIT_SHA1 := $(shell git rev-parse HEAD)
AppName := attacker

default: attacker

all: attacker reward

BUILD_FLAGS = -tags netgo -ldflags "\
	-X github.com/tsinghua-cel/attacker-service/versions.AppName=${AppName} \
	-X github.com/tsinghua-cel/attacker-service/versions.TagVersion=${VERSION} \
	-X 'github.com/tsinghua-cel/attacker-service/versions.BuildTime=`date`' \
	-X github.com/tsinghua-cel/attacker-service/versions.CommitSha1=${COMMIT_SHA1}  \
	-X 'github.com/tsinghua-cel/attacker-service/versions.GoVersion=`go version`' \
	-X 'github.com/tsinghua-cel/attacker-service/versions.GitBranch=`git symbolic-ref --short -q HEAD`' \
	"

attacker:
	go build $(BUILD_FLAGS) -o=${GOBIN}/$@ -gcflags "all=-N -l" ./cmd/attacker
	@echo "Done building."

reward:
	go build $(BUILD_FLAGS) -o=${GOBIN}/$@ -gcflags "all=-N -l" ./cmd/rewards
	@echo "Done building."

docs:
	@swag init -g ./openapi/server.go

clean:
	rm -fr build/*
docker:
	docker build -t attacker:${TAG} .
