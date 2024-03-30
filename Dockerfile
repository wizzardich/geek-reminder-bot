FROM golang:1.22-alpine as builder

ENV SRC_DIR=/go/src/github.com/wizzardich/geek-reminder-bot/

WORKDIR $SRC_DIR

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . $SRC_DIR

RUN go build -o /app/geek-reminder-bot

FROM golang:1.22-alpine 

LABEL org.opencontainers.image.source="https://github.com/wizzardich/geek-reminder-bot" \
      org.opencontainers.image.title="Geek Reminder Bot" \
      org.opencontainers.image.description="A Telegram Bot backend that serves as a Doodle scheduler" \
      org.opencontainers.image.authors="wizzardich" \
      org.opencontainers.image.licenses="MIT"

COPY --from=builder /app/geek-reminder-bot /app/geek-reminder-bot
WORKDIR /app

EXPOSE 8443

ENTRYPOINT ["./geek-reminder-bot"]