package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Initiator of the poll; its owner
type Initiator struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Notify   bool   `json:"notify"`
	TimeZone string `json:"timeZone"`
}

// A DateOption represents a pickable option in the poll
type DateOption struct {
	AllDay bool   `json:"allday"`
	Start  int64  `json:"start"`
	End    *[]int `json:"end"`
	ID     string `json:"id"`
}

// A PollRequest contains all information necessary to create a Doodle Poll
type PollRequest struct {
	Initiator       Initiator    `json:"initiator"`
	Options         []DateOption `json:"options"`
	Participants    []string     `json:"participants"`
	Comments        []string     `json:"comments"`
	Type            string       `json:"type"`
	Title           string       `json:"title"`
	Description     string       `json:"description"`
	PreferencesType string       `json:"preferencesType"`
	Hidden          bool         `json:"hidden"`
	RemindInvitees  bool         `json:"remindInvitees"`
	AskAddress      bool         `json:"askAddress"`
	AskEmail        bool         `json:"askEmail"`
	AskPhone        bool         `json:"askPhone"`
	Locale          string       `json:"locale"`
}

// A PollCreated object represents a simplified JSON response from Doodle API
type PollCreated struct {
	ID       string `json:"id"`
	AdminKey string `json:"adminKey"`
}

func composeOptions() *[]DateOption {
	options := make([]DateOption, 0, 7)

	now := time.Now()

	start := now
	for start.Weekday() != time.Monday {
		start = start.AddDate(0, 0, 1)
	}

	index := 0
	current := time.Date(start.Year(), start.Month(), start.Day(), 12, 0, 0, 0, time.UTC)
	for current == start || current.Weekday() != time.Monday {
		newUUID := uuid.NewV4()
		options = append(options, DateOption{true, current.Unix() * 1000, nil, newUUID.String()})
		index++
		current = current.AddDate(0, 0, 1)
	}

	return &options
}

func composeTitle() string {
	now := time.Now()

	start := now
	for start.Weekday() != time.Monday {
		start = start.AddDate(0, 0, 1)
	}

	end := start.AddDate(0, 0, 1)
	for end.Weekday() != time.Sunday {
		end = end.AddDate(0, 0, 1)
	}

	title := start.Format("Jan 02") + " - " + end.Format("Jan 02") + " Geek Availability"

	return title
}

func newPollRequest() *PollRequest {
	initiator := Initiator{"Your friendly bot", hostEmail, true, hostTimeZone}
	options := composeOptions()
	title := composeTitle()

	return &PollRequest{initiator, *options, []string{}, []string{}, "DATE", title, "", "YESNOIFNEEDBE", false, false, false, false, false, "en_US"}
}

func createPoll() (*PollCreated, error) {
	url := "https://doodle.com/api/v2.0/polls"

	pollRequest := newPollRequest()

	log.Printf("Preparing poll request %+v", pollRequest)

	jsonB, err := json.Marshal(pollRequest)

	if err != nil {
		log.Println("Error while marshalling poll request object.")
		log.Println(err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonB))

	if err != nil {
		log.Println("Error while creating a HTTP request object.")
		log.Println(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Error while performing an HTTP request object to %s: %+v.", url, req)
		log.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	var responseBody PollCreated
	err = json.NewDecoder(resp.Body).Decode(&responseBody)

	if err != nil {
		log.Printf("Error while decoding Doodle response object: %+v.", resp)
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Printf("Received response with status code set to %d.\n", resp.StatusCode)
		log.Printf("Response body: %+v", responseBody)
		return nil, fmt.Errorf("response from doodle with status %d", resp.StatusCode)
	}

	return &responseBody, nil
}
