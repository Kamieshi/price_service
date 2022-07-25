#!/bin/bash
export PATH=$(go env GOPATH)/bin:$PATH
protoc -I . ./protoc/priceService.proto --go_out=:.
protoc -I . ./protoc/priceService.proto --go-grpc_out=:.

