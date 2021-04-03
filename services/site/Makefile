.PHONY: static
static: 
	CGO_ENABLED=0 go build -a -ldflags '-s' -o site_service

.PHONY: docker
docker: static
	docker build -t kira/site -f ./docker/Dockerfile .