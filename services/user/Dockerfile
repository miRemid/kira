FROM alpine:3.7

ADD user_service /

RUN chmod 777 /user_service

RUN ls /

ENTRYPOINT [ "sh", "-c", "/user_service" ]