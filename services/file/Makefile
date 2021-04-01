.PHONY: build
build: 
	go build -o file_service main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build -t kira/file -f ./docker/Dockerfile .

.PHONY: static
static:
	CGO_ENABLED=0 go build -a -ldflags '-s' -o file_service