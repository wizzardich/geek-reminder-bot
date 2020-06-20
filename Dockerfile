FROM golang:1.14-alpine

LABEL org.opencontainers.image.source="https://github.com/wizzardich/geek-reminder-bot" \
      org.opencontainers.image.title="Geek Reminder Bot" \
      org.opencontainers.image.description="A Telegram Bot backend that serves as a Doodle scheduler" \
      org.opencontainers.image.authors="wizzardich" \
      org.opencontainers.image.licenses="MIT"

ENV SRC_DIR=/go/src/github.com/wizzardich/geek-reminder-bot/

WORKDIR /app
ADD . $SRC_DIR
RUN cd $SRC_DIR && \
    go build -o /app/geek-reminder-bot

EXPOSE 8443

ENTRYPOINT ["./geek-reminder-bot"]