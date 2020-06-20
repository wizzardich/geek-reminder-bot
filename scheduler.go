package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Scheduler contains current time till the ticker, the ticker, the responder bot API and the killer channel
type Scheduler struct {
	timer  *time.Timer
	ticker *time.Ticker
	bot    *tgbotapi.BotAPI
	kill   chan bool
}

var scheduleCache = make(map[int64]*Scheduler)
var timeout = 24 * time.Hour

func revoke(channelID int64) {
	if scheduler, ok := scheduleCache[channelID]; ok {
		go func() {
			scheduler.kill <- true
			log.Printf("[%d] -- succesfully unscheduled", channelID)
		}()
		msg := tgbotapi.NewMessage(channelID, "Отменили все расписания в этом канале, капитан!")
		scheduler.bot.Send(msg)
	}
}

func invoke(scheduler *Scheduler, channelID int64) {
	processChannel := make(chan bool)
	go func() { processChannel <- true }()
	select {
	case <-scheduler.kill:
		log.Printf("[%d] closing the scheduler", channelID)
		if scheduler.ticker != nil {
			scheduler.ticker.Stop()
		}
		scheduler.timer.Stop()
		delete(scheduleCache, channelID)
		return
	case <-scheduler.timer.C:
		scheduler.ticker = time.NewTicker(timeout)
	}

	for {
		select {
		case <-scheduler.kill:
			log.Printf("[%d] closing the scheduler", channelID)
			if scheduler.ticker != nil {
				scheduler.ticker.Stop()
			}
			scheduler.timer.Stop()
			delete(scheduleCache, channelID)
			return
		case <-processChannel:
			now := time.Now()
			if now.Weekday() == time.Sunday {
				log.Printf("[%d] Scheduling a new poll for %s.\n", channelID, now.String())
				scheduleNowDoodle(scheduler.bot, channelID)
			}
		case <-scheduler.ticker.C:
			now := time.Now()
			if now.Weekday() == time.Sunday {
				log.Printf("[%d] Scheduling a new poll for %s.\n", channelID, now.String())
				scheduleNowDoodle(scheduler.bot, channelID)
			}
		}
	}
}

func scheduleWeeklyDoodle(bot *tgbotapi.BotAPI, channelID int64) {
	now := time.Now()

	if _, ok := scheduleCache[channelID]; ok {
		msg := tgbotapi.NewMessage(channelID, "А мы уже, капитан!")
		bot.Send(msg)
		return
	}

	var nearest10AM time.Time
	if now.Hour() < 10 {
		nearest10AM = time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, time.UTC)
	} else {
		tomorrow := now.AddDate(0, 0, 1)
		nearest10AM = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, time.UTC)
	}
	delta := nearest10AM.Sub(time.Now())

	log.Printf("[%d] -- new channel scheduled it's first check in %s\n", channelID, delta.String())

	timer := time.NewTimer(delta)
	scheduler := Scheduler{timer, nil, bot, make(chan bool)}

	scheduleCache[channelID] = &scheduler

	go invoke(&scheduler, channelID)

	msg := tgbotapi.NewMessage(channelID, "Капитан, планируем планировать! Прям через "+delta.String()+" проверю иль не пора уж!")
	scheduler.bot.Send(msg)
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
