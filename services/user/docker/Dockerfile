FROM alpine:3.7

RUN mkdir /user

ADD user_service /user

RUN chmod 777 /user/user_service

WORKDIR /user

ENTRYPOINT [ "sh", "-c", "./user_service" ]