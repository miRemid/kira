GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: proto
proto:
	protoc --proto_path=. proto/*.proto --micro_out=pb --go_out=pb
	protoc-go-inject-tag -input=pb/file.pb.go

.PHONY: build
build: proto
	go build -o file_service main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t file_service:latest
