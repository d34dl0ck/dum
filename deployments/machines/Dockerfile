# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o /machines-service ./cmd/machines

EXPOSE 3000

CMD [ "/machines-service" ]