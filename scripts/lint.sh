#!/usr/bin/env bash

export GO111MODULE=on

go get -t -v ./...

golangci-lint run -v --timeout 2m -c .golangci.yml ./...