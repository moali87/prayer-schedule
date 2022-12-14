// schedule test for pcal which returns monthly prayer data based on customer input
package schedule_test

import (
	"fmt"
	"os"
	"testing"
	"time"
    psched "github.com/moali87/prayer-schedule"
)

func TestPrayerCalendarWithAPIKey(t *testing.T) {
	beverlyHillsTimeZone, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		errMsg := fmt.Errorf("Unable to load timezone data for America/Los_Angeles")
		panic(errMsg)
	}
	beverlyHillsTime := time.Date(2022, time.October, 22, 10, 10, 0, 0, beverlyHillsTimeZone)
	// customerInputWithAPIKey := &schedule.CustomerLocationInput{
	customerInputWithAPIKey := &psched.CustomerLocationInput{
		CountryCode: "USA",
		HEREAPIKey:  os.Getenv("HERE_API_KEY"),
		PostalCode:  "90210",
		CustTime:    beverlyHillsTime,
	}

	hereInput := new(psched.CustomerLocationInputWithHEREAPIKey)
	hereInput.CountryCode = customerInputWithAPIKey.CountryCode
	hereInput.PostalCode = customerInputWithAPIKey.PostalCode
	hereInput.HEREAPIKey = customerInputWithAPIKey.HEREAPIKey
	hereResp, hereAddressData, err := psched.HERECustomerLocation(hereInput)
	if err != nil {
		t.Errorf("HERE returned an error: %v", err)
	}
	t.Logf("HERE full response: %v", hereResp)
	t.Logf("HERE addresss data: %v", hereAddressData)
	t.Logf("HERE Longitude: %v", hereAddressData.Coordiantes.Lng)
	t.Logf("HERE Latitude: %v", hereAddressData.Coordiantes.Lat)

	monthlyData, err := customerInputWithAPIKey.PrayerCalendar()
	if err != nil {
		t.Errorf("error looking up customer data with api key: %v", err)
	}

	if monthlyData.Code != 200 {
		t.Errorf("customerInputWithAPIKey returned code is not 200: %v", err)
	}

	if len(monthlyData.Data) == 0 {
		t.Error("monthly data did not return any timings")
	}

	fmt.Printf("Some prayer data with API Key %v", monthlyData.Data[0].Timings.Asr)
}

func TestPrayerCalendarWithoutAPIKey(t *testing.T) {
	beverlyHillsTimeZone, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		errMsg := fmt.Errorf("Unable to load timezone data for America/Los_Angeles")
		panic(errMsg)
	}
	beverlyHillsTime := time.Date(2022, time.October, 22, 10, 10, 0, 0, beverlyHillsTimeZone)

	customerInputWithAPIKey := &psched.CustomerLocationInput{
		Coordinates: psched.PrayerCalendarInputCoordinates{
			Latitude:  34.1030,
			Longitude: -118.4105,
		},
		CustTime: beverlyHillsTime,
	}

	monthlyData, err := customerInputWithAPIKey.PrayerCalendar()
	if err != nil {
		t.Errorf("error looking up customer data with api key: %v", err)
	}

	if monthlyData.Code != 200 {
		t.Errorf("customerInputWithAPIKey returned code is not 200: %v", err)
	}

	if len(monthlyData.Data) == 0 {
		t.Error("monthly data did not return any timings")
	}

	fmt.Printf("Some prayer data without API Key: %v", monthlyData.Data[0].Timings.Asr)
}
