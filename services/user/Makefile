.PHONY: build
build: 
	go build -o user_service main.go

.PHONY: run
run:
	./user_service

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: wrk
wrk:
	wrk -c 100 -t 12 -d 30s --latency -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjI4MDE1MDIsIlVzZXJJRCI6ImMyb2MzajlkM2NsNWRpMnZuZWswIiwiUm9sZSI6Im5vcm1hbCJ9.qsUmdaVBy7jioHVynVVMNuSZeme-MLOxhiQ2nK_4W-c5Zn7ezLTr9xF24bV8-KpEitgoVYfcXMKQrX8v9ttVVB6y-ghaXlVsNKy1xqgZKqH9qfgdZOIKjC1qP2vok1hSyMsc2d-nEShSPf0Bl14Y_yIMEa6NhOk1PNId0lwGA_k" http://api.test.me/user/me

.PHONY:
test-image:
	wrk -c 100 -t 12 -d 30s --latency http://api.test.me/image?id=YJnyxmcuI&gray=true

.PHONY: docker
docker: static
	docker build -t kira/user -f ./docker/Dockerfile .

.PHONY: static
static:
	CGO_ENABLED=0 go build -a -ldflags '-s' -o user_service

