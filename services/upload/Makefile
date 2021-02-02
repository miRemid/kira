.PHONY: build
build: 
	go build -o upload_service main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t upload_service:latest
