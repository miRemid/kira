.PHONY: build
build: 
	go build -o file_service main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t file_service:latest
