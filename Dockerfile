FROM golang:1.11-alpine as builder
RUN apk update && apk upgrade && apk add --no-cache bash git
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GO11MODULE=ON go build .

FROM scratch 
COPY --from=builder /app/go-rest-api .
COPY --from=builder /app/.env .

EXPOSE 8081
ENTRYPOINT ["/go-rest-api"]"
