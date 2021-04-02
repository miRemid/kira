export MICRO_REGISTRY=etcd
export MICRO_REGISTRY_ADDRESS=127.0.0.1:2379
export MICRO_API_HANDLER=http
export MICRO_NAMESPACE=kira.micro.api
export CORS_ALLOWED_METHODS="HEAD,GET,POST,OPTIONS,PUT,DELETE"
export CORS_ALLOWED_HEADERS="Origin,Content-Type,Accept,Authorization,Content-Length,Accept-Encoding,X-CSRF-Token"
export CORS_ALLOWED_ORIGINS="*"
export CORS_ALLOWED_CREDENTIALS=true

.PHONY: build
build:
	go build -o gateway_service

.PHONY: docker
docker:
	docker build -t kira/gateway -f ./docker/Dockerfile .

.PHONY: static
static:
	CGO_ENABLED=0 go build -a -ldflags '-s' -o gateway_service

.PHONY: run
run: build
	MICRO_NAMESPACE=kira.micro.api ./gateway_service api