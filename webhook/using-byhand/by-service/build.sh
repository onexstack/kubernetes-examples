#!/bin/bash

: ${REGISTRY_PREFIX:? required}

export GO111MODULE=on
export GOPROXY=https://goproxy.cn
# build webhook
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o admission-webhook-example
# build docker image
docker build --no-cache -t ${REGISTRY_PREFIX}/admission-webhook-example:v1 .
rm -rf admission-webhook-example

docker push ${REGISTRY_PREFIX}/admission-webhook-example:v1
