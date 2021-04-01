FROM alpine:3.7 

ADD auth_service /
ADD pem /pem

RUN chmod 777 /auth_service

ENTRYPOINT [ "sh", "-c", "/auth_service" ]