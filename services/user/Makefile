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
	wrk -c 100 -t 5 -d 1m -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTQ0MjM1NzQsIlVzZXJJRCI6ImMwbzhvb3Y2OWtmNW0xYW10NDRnIiwiUm9sZSI6ImFkbWluIn0.yey4dA_s7Mi1_5rh574SOi2tfOZw7cDMkNiIctF3MeXBHrB3j65KAHoIL1bjvtJc04XFM7wyGptT7fnizjm-WNjqz7IvfkbmjaXy0ft6l3g-JHyBQrTRdlhkZICwiw_8kS5wAWHjQloEm7SukNXDIDMG92heyGnu13pnNZDvt30" http://api.test.me/v1/user/me

.PHONY: docker
docker:
	docker build -t kira/user -f ./docker/Dockerfile .

.PHONY: static
static:
	CGO_ENABLED=0 go build -a -ldflags '-s' -o user_service

