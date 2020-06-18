package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	token := os.Getenv("GO_TELEGRAM_TOKEN")

	if token == "" {
		log.Fatalln("Environment variable GO_TELEGRAM_TOKEN is not defined.")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://409f55f9ebfd.ngrok.io/" + bot.Token))

	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()

	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)

	go http.ListenAndServe(":8443", nil)

	for update := range updates {
		log.Printf("%+v\n", update)

		if update.ChannelPost == nil {
			continue
		}

		log.Printf("[%d] -- %+v", update.ChannelPost.Chat.ID, update.ChannelPost)

		msg := tgbotapi.NewMessage(update.ChannelPost.Chat.ID, "Ahoy, mateys!")
		msg.ReplyToMessageID = update.ChannelPost.MessageID

		bot.Send(msg)
	}
}
