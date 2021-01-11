FROM alpine:3.7

ADD bin/user_service /

RUN chmod 777 /user_service

ENTRYPOINT [ "sh", "-c", "/user_service" ]