FROM alpine:3.7

ADD upload_api_service /
RUN chmod 777 /upload_api_service

ENTRYPOINT [ "sh", "-c", "/upload_api_service" ]