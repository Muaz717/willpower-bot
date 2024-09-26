FROM golang:1.23.1-alpine3.20 AS builder

COPY . /github.com/Muaz717/willpower-bot/
WORKDIR /github.com/Muaz717/willpower-bot/

RUN go mod download

RUN GOOS=linux go build -o ./.bin/bot ./cmd/bot/main.go 

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/Muaz717/willpower-bot/.bin/bot .
COPY --from=builder /github.com/Muaz717/willpower-bot/config/config.yml config/
COPY --from=builder /github.com/Muaz717/willpower-bot/.env .

EXPOSE 80

CMD ["./bot"]
