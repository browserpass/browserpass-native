FROM golang:latest

VOLUME /src
WORKDIR /src

ENTRYPOINT ["make"]
