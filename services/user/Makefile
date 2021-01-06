GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: proto
proto:
	protoc --proto_path=. proto/*.proto --micro_out=pb --go_out=pb
	protoc-go-inject-tag -input=pb/user.pb.go

.PHONY: build
build: proto
	go build -o user_service main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t user_service:latest
