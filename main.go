package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var hostEmail string
var hostTimeZone string

const tokenEnv = "GO_TELEGRAM_TOKEN"
const emailEnv = "GO_DOODLE_EMAIL"
const timezEnv = "GO_DOODLE_TZ"
const localEnv = "GO_EXTERNAL_ADDRESS"

func main() {
	token := os.Getenv(tokenEnv)

	if token == "" {
		log.Fatalf("Environment variable %s is not defined.\n", tokenEnv)
	}

	hostEmail = os.Getenv(emailEnv)

	if hostEmail == "" {
		log.Fatalf("Environment variable %s is not defined.\n", emailEnv)
	}

	localURL := os.Getenv(localEnv)

	if hostEmail == "" {
		log.Fatalf("Environment variable %s is not defined.\n", localEnv)
	}

	hostTimeZone = os.Getenv(timezEnv)

	if hostTimeZone == "" {
		log.Printf("Environment variable %s is not defined, using UTC timezone.\n", timezEnv)
		hostTimeZone = "UTC"
	}

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Fatal(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(localURL + bot.Token))

	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()

	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Previous telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)

	go http.ListenAndServe(":8443", nil)

	for update := range updates {
		if update.ChannelPost == nil {
			continue
		}

		log.Printf("[%d] -- %s", update.ChannelPost.Chat.ID, update.ChannelPost.Text)

		msg := tgbotapi.NewMessage(update.ChannelPost.Chat.ID, "Ahoy, mateys!")
		msg.ReplyToMessageID = update.ChannelPost.MessageID

		if update.ChannelPost.Text == "GO" {
			created, err := createPoll()

			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Poll %s created with admin key %s", created.ID, created.AdminKey)
		}

		bot.Send(msg)
	}
}
