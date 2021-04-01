FROM alpine:3.7 

ADD file_service /

RUN chmod 777 /file_service

ENTRYPOINT [ "sh", "-c", "/file_service" ]