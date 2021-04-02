GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/api.proto=github.com/micro/go-micro/v2/api/proto
export MICRO_REGISTRY=etcd
export MICRO_REGISTRY_ADDRESS=127.0.0.1:2379
export MICRO_API_HANDLER=http
export MICRO_NAMESPACE=kira.micro.api

.PHONY: build
build:
	sh build.sh

.PHONY: proto
proto:
	protoc --proto_path=. proto/*.proto --micro_out=proto/pb --go_out=proto/pb
	protoc-go-inject-tag -input=proto/pb/auth.pb.go
	protoc-go-inject-tag -input=proto/pb/file.pb.go
	protoc-go-inject-tag -input=proto/pb/user.pb.go
	protoc-go-inject-tag -input=proto/pb/upload.pb.go

micro:
	micro web