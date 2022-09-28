// Package schedule returns customer monthly prayer schedule
package schedule

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// FiveDailyPrayers is all of the five prayers for the day in index +-1 as well as the time of prayer
type FiveDailyPrayers struct {
	Fajr    string `json:"Fajr"`
  Sunrise string `json:"Sunrise"`
	Dhuhr   string `json:"Dhuhr"`
	Asr     string `json:"Asr"`
	Maghrib string `json:"Maghrib"`
	Isha    string `json:"Isha"`
}

// PCalInput is the customer geolocation and prayer source method
type PCalInput struct {
	CustTime    time.Time
	Institution int     // Aladhan prayer data source method
	Latitude    float32 // Client latitude to use with aladhan
	Longitude   float32 // Client longitude to use with aladhan
}

// PCalOutput contains the prayer time of the month as well as the return code
type PCalOutput struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   []struct {
		Timings FiveDailyPrayers
	}
	Latitude  float32 // Client latitude to use with aladhan
	Longitude float32 // Client longitude to use with aladhan
}

/*
AladhanData returns the total monthly prayers of given month, coordinates, and zip from aladhan.
https://api.aladhan.com/v1/calendar?latitude=51.508515&longitude=-0.1254872&method=1&month=4&year=2017
*/
func AladhanData(input *PCalInput) (*PCalOutput, error) {
	// Use HERE API to get client coordinates
	reqURL := fmt.Sprintf(
		"https://api.aladhan.com/v1/calendar?latitude=%v&longitude=%v&method=%d&month=%d&year=%d",
		input.Latitude,
		input.Longitude,
		input.Institution,
		input.CustTime.Month(),
		input.CustTime.Year(),
	)

	resp := new(PCalOutput)
	req, err := http.Get(reqURL)

	if err != nil {
		fmt.Println(reqURL)
		fmt.Println(err)
	}

	defer func() {
		err := req.Body.Close()
		if err != nil {
			fmt.Println("Failed to close request from Aladhan")
			fmt.Println(err)
			panic(err)
		}
	}()

	json.NewDecoder(req.Body).Decode(resp)
	if resp.Code != 200 {
		fmt.Println("Aladhan response is not 200")
		fmt.Println(resp)
		panic("Err")
	}

	return resp, nil
}
