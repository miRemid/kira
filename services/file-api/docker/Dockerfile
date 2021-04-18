FROM alpine:3.7

ADD file_api_service /
ADD casbin /casbin
RUN chmod 777 /file_api_service

ENTRYPOINT [ "sh", "-c", "/file_api_service" ]