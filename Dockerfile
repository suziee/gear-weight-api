FROM golang:1.23.3-alpine AS base
WORKDIR /home/app

RUN apk update \
    && apk add --no-cache protobuf \
    && apk add build-base \
    && rm -rf /var/cache/apk/* \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && export PATH="$PATH:$(go env GOPATH)/bin" \
    && go env -w CGO_ENABLED=1 \
    && apk add git \
    && go install github.com/mattn/go-sqlite3@1.14.16 