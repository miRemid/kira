.PHONY: static
static:
	CGO_ENABLED=0 go build -a -ldflags '-s' -o upload_api_service

.PHONY: docker
docker:
	docker build -t kira/upload-api -f ./docker/Dockerfile .