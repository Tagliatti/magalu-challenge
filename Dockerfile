ARG DIR=/app

FROM golang:1.24-bullseye AS dev

ARG DIR
WORKDIR $DIR

RUN git config --global --add safe.directory $DIR &&  \
    go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

EXPOSE 8080
EXPOSE 2345

CMD ["air"]

FROM golang:1.24-bullseye AS build

ARG DIR
WORKDIR $DIR

COPY docker .

RUN git config --global --add safe.directory $DIR &&  \
    go build

FROM debian:bullseye-slim

ARG DIR
WORKDIR $DIR

COPY --from=build ${DIR}/challenge-bravo .

EXPOSE 8080

CMD ["./challenge-bravo"]
