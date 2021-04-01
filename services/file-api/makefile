build:
	go build -o file-api

static:
	CGO_ENABLED=0 go build -a -ldflags '-s' -o file_api_service

.PHONY: docker
docker:
	docker build -t kira/file-api -f ./docker/Dockerfile .