# syntax=docker/dockerfile:1

FROM golang:1.18.4

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /userservice

EXPOSE 8080

CMD [ "/userservice" ]