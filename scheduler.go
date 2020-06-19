package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func scheduleWeeklyDoodle(bot *tgbotapi.BotAPI, channelID int64) {

}

func scheduleNowDoodle(bot *tgbotapi.BotAPI, channelID int64) {
	created, err := createPoll()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Poll %s created with admin key %s", created.ID, created.AdminKey)

	msg := tgbotapi.NewMessage(channelID, "Ахой, гики! Еженедельный дудл подвезли! #doodle\nhttps://doodle.com/poll/"+created.ID)
	bot.Send(msg)
}
