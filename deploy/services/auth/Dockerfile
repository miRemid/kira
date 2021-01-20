FROM alpine:3.7 

ADD bin/auth_service /
ADD bin/pem /pem

RUN ls && chmod 777 /auth_service

ENTRYPOINT [ "sh", "-c", "/auth_service" ]