FROM alpine:3.7

ADD upload_service /

RUN chmod 777 /upload_service

ENTRYPOINT [ "sh", "-c", "/upload_service" ]