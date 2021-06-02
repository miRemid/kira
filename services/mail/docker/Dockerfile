FROM alpine:3.7

ADD mail_service /
RUN chmod 777 /mail_service

ENTRYPOINT [ "sh", "-c", "/mail_service" ]