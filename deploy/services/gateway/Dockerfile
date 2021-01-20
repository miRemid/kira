FROM alpine:3.7

ADD bin/gateway_service /
ADD bin/casbin /casbin

RUN chmod 777 /gateway_service

ENTRYPOINT [ "sh", "-c", "/gateway_service api --address=0.0.0.0:3000 --handler=web"]