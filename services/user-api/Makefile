.PHONY: static
static: 
	CGO_ENABLED=0 go build -a -ldflags '-s' -o user_api_service

.PHONY: docker
docker: static
	docker build -t kira/user-api -f ./docker/Dockerfile .