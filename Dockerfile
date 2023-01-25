FROM golang:1.20rc1-alpine3.17 AS build

WORKDIR /app

RUN apk add --update alpine-sdk

ADD go.mod .

COPY . .
RUN go get ./...
RUN go mod download

RUN go build -o main

FROM alpine

WORKDIR /app
COPY resources/config.yaml .

COPY resources/secrets ./secrets
COPY ../Worker/templates ./templates

COPY --from=build /app/main main
CMD ["./main", "./config.yaml"]