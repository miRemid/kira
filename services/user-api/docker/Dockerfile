FROM alpine:3.7

RUN mkdir /user

ADD user_api_service /user
ADD casbin /user/casbin

RUN chmod 777 /user/user_api_service

WORKDIR /user

ENTRYPOINT [ "sh", "-c", "./user_api_service" ]