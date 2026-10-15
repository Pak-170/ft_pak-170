# syntax=docker/dockerfile:1

FROM golang:1.27.1-alpine

ENV GOTOOLCHAIN=auto

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/docker-entrypoint.sh"]