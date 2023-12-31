FROM golang:1.15-alpine3.12 AS builder

COPY . /github.com/MoJIoToK/Go_Telebot_pocket/
WORKDIR /github.com/MoJIoToK/Go_Telebot_pocket/

RUN go mod download
RUN go build -o /bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/MoJIoToK/Go_Telebot_pocket/bin/bot .
COPY --from=0 /github.com/MoJIoToK/Go_Telebot_pocket/configs configs/

EXPOSE 80

CMD ["./bot"]