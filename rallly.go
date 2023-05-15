package main

import "time"

type RalllyPollInitiator struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RalllyPollCreate struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Options     []RalllyDateOption  `json:"options"`
	User        RalllyPollInitiator `json:"user"`
	Timezone    string              `json:"timezone"`
}

type RalllyDateOption struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate,omitempty"`
}

type RalllyPollCreateRequest struct {
	Root struct {
		JSON RalllyPollCreate `json:"json"`
	} `json:"0"`
}

type RalllyPollCreated struct {
	ID    string `json:"id"`
	URLID string `json:"urlId"`
}

type RalllyPollCreatedResponse struct {
	Result struct {
		Data struct {
			Json RalllyPollCreated `json:"json"`
		} `json:"data"`
	} `json:"result"`
}

func composeRalllyOptions() []RalllyDateOption {
	options := make([]RalllyDateOption, 0, 7)

	now := time.Now()

	start := now
	for start.Weekday() != time.Monday {
		start = start.AddDate(0, 0, -1)
	}

	end := start.AddDate(0, 0, 7)

	for start.Before(end) {
		options = append(options, RalllyDateOption{
			StartDate: start.Format("2006-01-02"),
		})
		start = start.AddDate(0, 0, 1)
	}

	return options
}

func newRalllyPollRequest() *RalllyPollCreateRequest {
	return &RalllyPollCreateRequest{
		Root: struct {
			JSON RalllyPollCreate `json:"json"`
		}{
			JSON: RalllyPollCreate{
				Title:       composeTitle(),
				Description: "Whelp, who could've thought that this would be so easy?",
				Options:     composeRalllyOptions(),
				User: RalllyPollInitiator{
					Name:  "Your friendly bot",
					Email: hostEmail,
				},
				Timezone: hostTimeZone,
			},
		},
	}
}

func createRalllyPoll() (*RalllyPollCreated, error) {
	var pollCreated RalllyPollCreatedResponse

	err := postJSON("https://schedule.smugglersden.org/api/trpc/polls.create", newRalllyPollRequest(), &pollCreated)
	if err != nil {
		return nil, err
	}

	return &pollCreated.Result.Data.Json, nil
}
