GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/api.proto=github.com/micro/go-micro/v2/api/proto


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
