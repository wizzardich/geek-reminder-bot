package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	End    string `json:"end"`
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

func composeOptions(start int64) []DateOption {
	options := make([]DateOption, 7)

	//TODO: define options truly

	return options
}

func newPollRequest(title string, start int64) PollRequest {
	initiator := Initiator{"Your friendly bot", hostEmail, true, hostTimeZone}
	options := composeOptions(start)

	return PollRequest{initiator, options, []string{}, []string{}, "DATE", title, "", "YESNOIFNEEDBE", false, false, false, false, false, "en_US"}
}

func createPoll(title string, start int64) error {
	url := "https://doodle.com/api/v2.0/polls"

	pollRequest := newPollRequest(title, start)

	log.Printf("Preparing poll request %+v", pollRequest)

	jsonB, err := json.Marshal(pollRequest)

	if err != nil {
		log.Println("Error while marshalling poll request object.")
		log.Println(err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonB))

	if err != nil {
		log.Println("Error while creating a HTTP request object.")
		log.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Error while performing an HTTP request object to %s: %+v.", url, req)
		log.Println(err)
		return err
	}

	defer resp.Body.Close()

	var responseBody PollCreated
	err = json.NewDecoder(resp.Body).Decode(&responseBody)

	if err != nil {
		log.Printf("Error while decoding Doodle response object: %+v.", resp)
		log.Println(err)
		return err
	}

	if resp.StatusCode != 200 {
		log.Printf("Received response with status code set to %d.\n", resp.StatusCode)
		log.Printf("Response body: %+v", responseBody)
		return fmt.Errorf("response from doodle with status %d", resp.StatusCode)
	}

	return nil
}
