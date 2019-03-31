FROM alpine:latest

ENV BINARY_NAME gorestapi

WORKDIR /app

COPY $BINARY_NAME .
COPY .env .

EXPOSE 8081
ENTRYPOINT [ "/app/gorestapi" ]
