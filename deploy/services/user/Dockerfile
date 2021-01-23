FROM alpine:3.7

ADD bin/user_service /
ADD bin/casbin /casbin

RUN chmod 777 /user_service

ENTRYPOINT [ "sh", "-c", "/user_service" ]