package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var hostEmail string
var hostTimeZone string

const tokenEnv = "GO_TELEGRAM_TOKEN"
const emailEnv = "GO_DOODLE_EMAIL"
const timezEnv = "GO_DOODLE_TZ"
const localEnv = "GO_EXTERNAL_ADDRESS"

const scheduleCommand = "/schedule"
const scheduleNowCommand = "/schedule now"
const scheduleWeeklyCommand = "/schedule weekly"
const unscheduleCommand = "/unschedule"

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
		switch {
		case update.ChannelPost == nil:
			continue
		case !update.ChannelPost.IsCommand():
			continue
		case strings.HasPrefix(update.ChannelPost.Text, scheduleNowCommand):
			log.Printf("[%d] -- %s", update.ChannelPost.Chat.ID, update.ChannelPost.Text)
			scheduleNowDoodle(bot, update.ChannelPost.Chat.ID)
		case strings.HasPrefix(update.ChannelPost.Text, scheduleWeeklyCommand):
			log.Printf("[%d] -- %s", update.ChannelPost.Chat.ID, update.ChannelPost.Text)
			scheduleWeeklyDoodle(bot, update.ChannelPost.Chat.ID)
		case strings.HasPrefix(update.ChannelPost.Text, unscheduleCommand):
			log.Printf("[%d] -- %s", update.ChannelPost.Chat.ID, update.ChannelPost.Text)
			revoke(update.ChannelPost.Chat.ID)
		default:
			log.Printf("[WARNING][%d] -- %s", update.ChannelPost.Chat.ID, update.ChannelPost.Text)
			continue
		}
	}
}
