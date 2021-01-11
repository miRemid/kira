FROM alpine:3.7

ADD bin/gateway_service /

RUN chmod 777 /gateway_service

ENTRYPOINT [ "sh", "-c", "/gateway_service" ]