FROM alpine:3.7

ADD bin/upload_service /

RUN chmod 777 /upload_service

ENTRYPOINT [ "sh", "-c", "/upload_service" ]