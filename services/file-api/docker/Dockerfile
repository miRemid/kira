FROM alpine:3.7

ADD file_api_service /
RUN chmod 777 /file_api_service

ENTRYPOINT [ "sh", "-c", "/file_api_service" ]