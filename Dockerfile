# Copyright (c) 2017 Alex Pliutau

# build
FROM golang:alpine as builder
RUN apk add --no-cache git
ADD . /go/src/github.com/wizeline/slack-stranger-bot
WORKDIR /go/src/github.com/wizeline/slack-stranger-bot
RUN go get github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go install

# binary only
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/bin/slack-stranger-bot .
ENTRYPOINT ["./slack-stranger-bot"]
