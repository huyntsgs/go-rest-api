FROM golang:1.11-alpine
RUN apk update && apk upgrade && apk add --no-cache bash git
RUN go get github.com/huyntsgs/go-rest-api

ENV SRC=/go/src/github.com/huyntsgs/go-rest-api
RUN cd ${SRC} $$ CGO_ENABLED=0 go build
WORKDIR ${SRC}
RUN cp .env /go/bin

EXPOSE 8081
CMD go-rest-api
