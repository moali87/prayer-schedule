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

type PCalOutputs struct {
	CurrentMonthCalendar  PCalOutput
	PreviousMonthCalendar PCalOutput
	NextMonthCalendar     PCalOutput
}

func BeginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func EndOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}

func aladhanReq(reqURL <-chan string, pcalOutput chan <-*PCalOutput) {

	resp := new(PCalOutput)
	req, err := http.Get(<-reqURL)

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

    pcalOutput <- resp

    close(pcalOutput)
}

/*
aladhanData returns the total monthly prayers of given month, coordinates, and zip from aladhan.
https://api.aladhan.com/v1/calendar?latitude=51.508515&longitude=-0.1254872&method=1&month=4&year=2017
*/
func AladhanData(input *PCalInput) (*PCalOutputs, error) {
	outputs := new(PCalOutputs)
	// Check to build if previous and next month calendar will be required
	var prevReqURL string
	var nextReqURL string
	if input.CustTime.Day() == BeginningOfMonth(input.CustTime).Day() {
		prevMonthDate := input.CustTime.AddDate(input.CustTime.Year(), input.CustTime.Day(), -1)
		prevReqURL = fmt.Sprintf(
			"https://api.aladhan.com/v1/calendar?latitude=%v&longitude=%v&method=%d&month=%d&year=%d",
			input.Latitude,
			input.Longitude,
			input.Institution,
			prevMonthDate.Month(),
			prevMonthDate.Year(),
		)
	}

	if input.CustTime.Day() == EndOfMonth(input.CustTime).Day() {
		nextMonthDate := input.CustTime.AddDate(input.CustTime.Year(), input.CustTime.Day(), 1)
		nextReqURL = fmt.Sprintf(
			"https://api.aladhan.com/v1/calendar?latitude=%v&longitude=%v&method=%d&month=%d&year=%d",
			input.Latitude,
			input.Longitude,
			input.Institution,
			nextMonthDate.Month(),
			nextMonthDate.Year(),
		)
	}

	// Use HERE API to get client coordinates
	reqURL := fmt.Sprintf(
		"https://api.aladhan.com/v1/calendar?latitude=%v&longitude=%v&method=%d&month=%d&year=%d",
		input.Latitude,
		input.Longitude,
		input.Institution,
		input.CustTime.Month(),
		input.CustTime.Year(),
	)

    var urlChan = make(chan string)

	if prevReqURL != "" {
        monthOutputChan := make(chan *PCalOutput)
		go aladhanReq(urlChan, monthOutputChan)
        urlChan <- prevReqURL
        fmt.Println("I've reached previousURL")
        monthOutput, ok := <-monthOutputChan
    if ok == false {
      panic("previous month output goroutine failed")
    }
    // close(monthOutputChan)
    outputs.PreviousMonthCalendar = *monthOutput
	}

	if nextReqURL != "" {
        monthOutputChan := make(chan *PCalOutput)
        go aladhanReq(urlChan, monthOutputChan)
        urlChan <- nextReqURL
        fmt.Println("I've reached nextURL")
        previousMonthOutput, ok := <-monthOutputChan
    if ok == false {
      panic("next month output goroutine failed")
    }
    // close(monthOutputChan)
    outputs.NextMonthCalendar = *previousMonthOutput
	}

    monthOutputChan := make(chan *PCalOutput)
    go aladhanReq(urlChan, monthOutputChan)
    urlChan <- reqURL
    fmt.Println("I've reached this point")
    monthOutput, ok := <-monthOutputChan
    if ok == false {
        panic("current month output goroutine failed")
    }

    // close(monthOutputChan)
	outputs.CurrentMonthCalendar = *monthOutput

	return outputs, nil
}

