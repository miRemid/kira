
.PHONY: build
build: 
	go build -o auth_service main.go

.PHONY: run
run: build
	MICRO_REGISTER=etcd MICRO_REGISTER_ADDRESS=127.0.0.1:2379 ./auth_service

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker: static
	docker build -t kira/auth -f ./docker/Dockerfile .

.PHONY: static	
static:
	CGO_ENABLED=0 go build -a -ldflags '-s' -o auth_service
