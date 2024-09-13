# build
FROM golang:1.23-alpine AS builder
RUN apk add build-base
WORKDIR /root
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o bot main.go

# binary only
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /root/bot .
ENTRYPOINT ["./bot"]
