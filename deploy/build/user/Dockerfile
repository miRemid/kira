FROM alpine:3.7

RUN mkdir /user

ADD bin/user_service /user
ADD bin/casbin /user/casbin

RUN chmod 777 /user/user_service

WORKDIR /user

ENTRYPOINT [ "sh", "-c", "./user_service" ]