# Copyright (c) 2017 Alex Pliutau

FROM golang:1.9

ENV PKG_PATH /go/src/github.com/wizeline/slack-stranger-bot

ADD . $PKG_PATH
WORKDIR $PKG_PATH

RUN curl https://glide.sh/get | sh
RUN glide install
RUN go install

ENTRYPOINT ["/go/bin/slack-stranger-bot"]
