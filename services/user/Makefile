.PHONY: build
build: 
	go build -o user_service main.go

.PHONY: run
run:
	./user_service

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t user_service:latest
