ARG DIR=/app
ARG MAIN

FROM golang:1.24-bullseye AS dev

ARG DIR
WORKDIR $DIR

RUN git config --global --add safe.directory $DIR &&  \
    go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

#EXPOSE 8080
EXPOSE 2345
EXPOSE 2346

FROM golang:1.24-bullseye AS build

ARG DIR
WORKDIR $DIR

COPY . .

RUN git config --global --add safe.directory $DIR &&  \
    go build -o main ./cmd/api/main.go

FROM debian:bullseye-slim

ARG DIR
WORKDIR $DIR

COPY --from=build ${DIR}/${MAIN} .

EXPOSE 8080

CMD ["./main"]
