#!/bin/bash
export PATH=$(go env GOPATH)/bin:$PATH
protoc -I . ./priceService.proto --go_out=:.
protoc -I . ./priceService.proto --go-grpc_out=:.

