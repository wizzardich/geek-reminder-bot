package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

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

func postJSON[IN, OUT any](url string, in IN, out OUT) error {
	jsonB, err := json.Marshal(in)

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

	err = json.NewDecoder(resp.Body).Decode(&out)

	if err != nil {
		log.Printf("Error while decoding response object: %+v.", resp)
		log.Println(err)
		return err
	}

	if resp.StatusCode != 200 {
		log.Printf("Received response with status code set to %d.\n", resp.StatusCode)
		log.Printf("Response body: %+v", out)
		return fmt.Errorf("response from rallly with status %d", resp.StatusCode)
	}

	return nil
}
