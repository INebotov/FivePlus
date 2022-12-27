FROM golang:1.14.4-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o main

FROM alpine

WORKDIR /app
COPY ./config.yaml .
COPY --from=build /app/main main
CMD ["./main", "./config.yaml"]