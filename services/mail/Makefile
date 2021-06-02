.PHONY: build
build: 
	go build -o mail_service main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker: static
	docker build -t kira/mail -f ./docker/Dockerfile .

.PHONY: static
static:
	CGO_ENABLED=0 go build -a -ldflags '-s' -o mail_service
