FROM alpine:3.7

ADD bin/site_service /

RUN chmod 777 /site_service

ENTRYPOINT [ "sh", "-c", "/site_service" ]