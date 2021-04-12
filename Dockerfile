FROM alpine:latest

ADD admission-webhook-cml-user /admission-webhook-cml-user
ENTRYPOINT ["./admission-webhook-cml-user"]