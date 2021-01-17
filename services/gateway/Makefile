GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: gorun
gorun:
	MICRO_REGISTRY=etcd MICRO_REGISTRY_ADDRESS=127.0.0.1:2379 go run main.go api --handler=web --address=0.0.0.0:3000

.PHONY: build
build: 
	go build -o gateway_service main.go

.PHONY: run
run: build
	MICRO_REGISTRY=etcd MICRO_REGISTRY_ADDRESS=127.0.0.1:2379 ./gateway_service api -handler=web --address=0.0.0.0:3000

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t gateway_service:latest
