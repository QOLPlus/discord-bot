# stage: builder
FROM golang:1.15-alpine as builder
ARG SSH_PRIV_KEY
WORKDIR /usr/src/app
COPY . .

RUN apk update && apk upgrade && apk add --no-cache bash git openssh \
    && git config --global url."ssh://git@github.com".insteadOf "https://github.com"
RUN mkdir -p /root/.ssh && \
    chmod 0700 /root/.ssh && \
    ssh-keyscan github.com > /root/.ssh/known_hosts && \
    echo "$SSH_PRIV_KEY" > /root/.ssh/id_rsa && \
    chmod 600 /root/.ssh/id_rsa
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPRIVATE=github.com/QOLPlus go get github.com/QOLPlus/core
RUN rm -rf /root/.ssh/id_rsa
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o discord-bot .

# stage: real image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/src/app/discord-bot ./discord-bot

CMD ["./discord-bot"]
