package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func restoreSchedule(bot *tgbotapi.BotAPI, registered *[]ChannelRecord) {
	for _, channel := range *registered {
		scheduler, delta := produceScheduler(bot)

		scheduleCache[channel.ChannelID] = scheduler

		go invoke(scheduler, channel.ChannelID)

		log.Printf("[%d] -- channel schedule restored, first check in %s\n", channel.ChannelID, (*delta).String())
	}
}

func produceScheduler(bot *tgbotapi.BotAPI) (*Scheduler, *time.Duration) {
	now := time.Now()

	var nearest10AM time.Time
	if now.Hour() < 10 {
		nearest10AM = time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, time.UTC)
	} else {
		tomorrow := now.AddDate(0, 0, 1)
		nearest10AM = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, time.UTC)
	}
	delta := time.Until(nearest10AM)

	log.Printf("new channel scheduled its first check in %s\n", delta.String())

	timer := time.NewTimer(delta)

	scheduler := Scheduler{timer, nil, bot, make(chan bool)}
	return &scheduler, &delta
}

func revoke(channelID int64) {
	if scheduler, ok := scheduleCache[channelID]; ok {
		go func() {
			scheduler.kill <- true
			log.Printf("[%d] -- succesfully unscheduled", channelID)
		}()
		msg := "Отменили все расписания в этом канале, капитан!"
		sendMsg(scheduler.bot, channelID, msg)
	}
}

func invoke(scheduler *Scheduler, channelID int64) {
	processChannel := make(chan bool, 1)
	processChannel <- true

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
				scheduleNowRallly(scheduler.bot, channelID)
			}
		case <-scheduler.ticker.C:
			now := time.Now()
			if now.Weekday() == time.Sunday {
				log.Printf("[%d] Scheduling a new poll for %s.\n", channelID, now.String())
				scheduleNowRallly(scheduler.bot, channelID)
			}
		}
	}
}

func scheduleNowRallly(bot *tgbotapi.BotAPI, channelID int64) {
	created, err := createRalllyPoll()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Rallly %s created", created.ID)

	msg := "Ахой, гики! Еженедельный ралли подвезли! #schedule\nhttps://" + ralllyEndpoint + "/admin/" + created.URLID
	sendMsg(bot, channelID, msg)
}

func scheduleWeeklyRallly(bot *tgbotapi.BotAPI, channelID int64) {
	if _, ok := scheduleCache[channelID]; ok {
		sendMsg(bot, channelID, "А мы уже, капитан!")
		return
	}

	scheduler, delta := produceScheduler(bot)

	scheduleCache[channelID] = scheduler

	go invoke(scheduler, channelID)

	msgText := "Капитан, планируем планировать! Прям через " + (*delta).String() + " проверю иль не пора уж!"
	sendMsg(scheduler.bot, channelID, msgText)
}

func sendMsg(bot *tgbotapi.BotAPI, channelID int64, msg string) {
	message := tgbotapi.NewMessage(channelID, msg)
	if _, err := bot.Send(message); err != nil {
		log.Printf("[%d] -- error while sending message: %s", channelID, err.Error())
	}
}
