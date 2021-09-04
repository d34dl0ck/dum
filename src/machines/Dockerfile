# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app/contracts/machines

COPY ./contracts/machines/ ./
RUN go mod download

WORKDIR /app/machines

COPY ./ ./
RUN go mod download

RUN go build -o /machines-service

EXPOSE 3000

CMD [ "/machines-service" ]