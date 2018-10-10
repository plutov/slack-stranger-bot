# Copyright (c) 2017 Alex Pliutau

# build
FROM golang:alpine as builder
RUN apk add --no-cache git gcc
ADD . /go/src/github.com/plutov/slack-stranger-bot
WORKDIR /go/src/github.com/plutov/slack-stranger-bot
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go install

# binary only
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/bin/slack-stranger-bot .
ENTRYPOINT ["./slack-stranger-bot"]
