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
	wrk -c 100 -t 10 -d 1m -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTMxMTc2MzksIlVzZXJJRCI6ImMwY2dpaXY2OWtmMGtkOHNncHQwIiwiUm9sZSI6ImFkbWluIn0.BNjCvcLDfXaY7RIDSI4ysv4UViW2_ZFGK-Rvx6y-IyXZFK0eQaQzH1RqWj70dwQ8WehTZFLDfWGPp_GH_JmfqawrGZr2Z-ng5mNpGXZritXA9kICI6iZTo641NkUfxuExGa8xPKMNHnhvq09p82-sAH5NtoV_YCYwtl_7TebK0Y" http://127.0.0.1:5002/v1/user/me

.PHONY: docker
docker:
	docker build . -t user_service:latest
