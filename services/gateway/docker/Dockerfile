FROM alpine:3.7

ADD gateway_service /
ADD casbin /

RUN chmod 777 /gateway_service

ENTRYPOINT [ "sh", "-c", "./gateway_service api"]