FROM golang:latest

ENV APP_PATH="$GOPATH/src/github.com/browserpass/browserpass-native"

RUN go get -u github.com/golang/dep/cmd/dep && \
    mkdir -p $APP_PATH && \
    ln -s $APP_PATH /

WORKDIR $APP_PATH

ENTRYPOINT ["make"]
