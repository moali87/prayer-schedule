package schedule

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// Tests returned data from Aladhan
func TestAladhanData(t *testing.T) {
	monthlyDataInput := &PCalInput{}

	/*
	  Test aladhan returned data from 90210 coordinates
	*/
	beverlyHillsTimeZone, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		errMsg := fmt.Errorf("Unable to load timezone data for America/Los_Angeles")
		panic(errMsg)
	}
	beverlyHillsTime := time.Date(2022, time.October, 22, 10, 10, 0, 0, beverlyHillsTimeZone)
	monthlyDataInput.Latitude = 34.1030
	monthlyDataInput.Longitude = 118.4105
	monthlyDataInput.Institution = 1
	monthlyDataInput.CustTime = beverlyHillsTime

	monthlyPrayerData, err := aladhanData(monthlyDataInput)
	if err != nil {
		aladhanReqFail := fmt.Errorf("Error: Call to Aladhan failed %s", err)
		t.Errorf(aladhanReqFail.Error())
	}

	if monthlyPrayerData.currentMonthCalendar.Code != 200 {
		t.Errorf("Error: Failed to retrieve data from aladhan")
	}
}

func TestAladhanDataEndOfMonth(t *testing.T) {
	monthlyDataInput := &PCalInput{}

	/*
	  Test aladhan returned data from 90210 coordinates
	*/
	beverlyHillsTimeZone, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		errMsg := fmt.Errorf("Unable to load timezone data for America/Los_Angeles")
		panic(errMsg)
	}
	beverlyHillsTime := time.Date(2021, time.December, 31, 10, 10, 0, 0, beverlyHillsTimeZone)
	monthlyDataInput.Latitude = 34.1030
	monthlyDataInput.Longitude = 118.4105
	monthlyDataInput.Institution = 1
	monthlyDataInput.CustTime = beverlyHillsTime

	monthlyPrayerData, err := aladhanData(monthlyDataInput)
	if err != nil {
		aladhanReqFail := fmt.Errorf("Error: Call to Aladhan failed %s", err)
		t.Errorf(aladhanReqFail.Error())
	}

	if monthlyPrayerData.currentMonthCalendar.Code != 200 {
		t.Errorf("Error: Failed to retrieve data from aladhan")
	}

	if monthlyPrayerData.nextMonthCalendar.Code != 200 {
		t.Errorf("Error: Next month calendar data was not retrieved from aladhan")
	}
}

func TestAladhanDataBeginningOfMonth(t *testing.T) {
	monthlyDataInput := &PCalInput{}

	/*
	  Test aladhan returned data from 90210 coordinates
	*/
	beverlyHillsTimeZone, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		errMsg := fmt.Errorf("Unable to load timezone data for America/Los_Angeles")
		panic(errMsg)
	}
	beverlyHillsTime := time.Date(2021, time.January, 1, 10, 10, 0, 0, beverlyHillsTimeZone)
	monthlyDataInput.Latitude = 34.1030
	monthlyDataInput.Longitude = 118.4105
	monthlyDataInput.Institution = 1
	monthlyDataInput.CustTime = beverlyHillsTime

	monthlyPrayerData, err := aladhanData(monthlyDataInput)
	if err != nil {
		aladhanReqFail := fmt.Errorf("Error: Call to Aladhan failed %s", err)
		t.Errorf(aladhanReqFail.Error())
	}

	if monthlyPrayerData.currentMonthCalendar.Code != 200 {
		t.Errorf("Error: Failed to retrieve data from aladhan")
	}

	if monthlyPrayerData.previousMonthCalendar.Code != 200 {
		t.Errorf("Error: Next month calendar data was not retrieved from aladhan")
	}
}

func TestHEREAladhan(t *testing.T) {
	customerLocationInput := &CustomerLocationInputWithHEREAPIKey{
		HEREAPIKey:  os.Getenv("HERE_API_KEY"),
		CountryCode: "USA",
		PostalCode:  "90210",
	}

	custLocRet, custCityLocRet, err := HERECustomerLocation(customerLocationInput)
	if err != nil {
		t.Errorf("Customer location lookup failed %s", err.Error())
	}

	beverlyHillsTimeZone, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		errMsg := fmt.Errorf("Unable to load timezone data for America/Los_Angeles")
		panic(errMsg)
	}
	beverlyHillsTime := time.Date(2022, time.October, 22, 10, 10, 0, 0, beverlyHillsTimeZone)
	custPcalInput := &PCalInput{}
	custPcalInput.Latitude = float32(custLocRet.Items[0].Position.Lat)
	custPcalInput.Longitude = float32(custLocRet.Items[0].Position.Lng)
	custPcalInput.Institution = 1
	custPcalInput.CustTime = beverlyHillsTime

	monthlyPrayerData, err := aladhanData(custPcalInput)
	if err != nil {
		aladhanReqFail := fmt.Errorf("Error: Call to Aladhan failed %s", err)
		t.Errorf(aladhanReqFail.Error())
	}

	if custCityLocRet.PostalCode != customerLocationInput.PostalCode {
		t.Error("Error: Failed to match postal code to HERE address output")
	}

	t.Log(custCityLocRet.Coordiantes.Lat)
	t.Log(custCityLocRet.Coordiantes.Lng)
	t.Log(monthlyPrayerData.currentMonthCalendar.Data[0].Timings)

	if monthlyPrayerData.currentMonthCalendar.Code != 200 {
		t.Error("Error: Failed to retrieve data from aladhan")
	}
}
