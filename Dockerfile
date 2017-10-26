# Copyright (c) 2017 Alex Pliutau

FROM golang:1.9

ENV PKG_PATH /go/src/github.com/wizeline/slack-stranger-bot

ADD . $PKG_PATH
WORKDIR $PKG_PATH

# Install dependencies
RUN go get github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go install

ENTRYPOINT ["/go/bin/slack-stranger-bot"]
