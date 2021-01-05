# stage: builder
FROM golang:1.15-alpine as builder
WORKDIR /usr/src/app
COPY . .

RUN apk update && apk upgrade && apk add --no-cache bash git openssl
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPRIVATE=github.com/QOLPlus go get github.com/QOLPlus/core
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o discord-bot .

# stage: real image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/src/app/discord-bot ./discord-bot

CMD ["./discord-bot"]
