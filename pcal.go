// Package schedule returns cutomer monthly prayer data

package schedule

import (
	"fmt"
	"os"
	"time"
)

type CustomerLocationInput struct {
	Coordinates PrayerCalendarInputCoordinates // Only required if HEREAPIKey is not filled
	CountryCode string
	CustTime    time.Time
	HEREAPIKey  string // Only required if Coordinates is not filled
	Institution int
	PostalCode  string // Only required if Coordiantes is not filled
}

type PrayerCalendarInputCoordinates struct {
	Latitude  float32
	Longitude float32
}

func NewPrayerCalendarWithCoordinates(
	customerTime time.Time,
	institution int,
	latitude float32,
	longitude float32) (*CustomerLocationInput, error) {
	return &CustomerLocationInput{
		Coordinates: PrayerCalendarInputCoordinates{
			Latitude:  latitude,
			Longitude: longitude,
		},
		CustTime:    customerTime,
		Institution: institution,
	}, nil
}

func NewPrayerCalendarWithoutCoordiantes(
	countryCode string,
	customerTime time.Time,
	institution int,
	hereAPIKey string,
	postalCode string) (*CustomerLocationInput, error) {
	return &CustomerLocationInput{
		CountryCode: countryCode,
		CustTime:    customerTime,
		HEREAPIKey:  hereAPIKey,
		Institution: institution,
		PostalCode:  postalCode,
	}, nil
}

/*
PrayerCalendar returns customer monthly prayer data with or without customer providing coordinates.
if customer does not provide coordiantes, they must provide a HERE API Key
*/
func (c *CustomerLocationInput) PrayerCalendar() (*PCalOutput, error) {
	lookupMethod, err := c.checkCustomerInput()
	if err != nil {
       fmt.Println(err) 
       return nil, err
	}

	hereLookup := new(CustomerLocationInputWithHEREAPIKey)
	monthlyPrayerData := new(PCalInput)

	monthlyPrayerData.CustTime = c.CustTime
	monthlyPrayerData.Institution = c.Institution
	hereLookup.CountryCode = c.CountryCode

	if lookupMethod != "Coordinates" && lookupMethod != "APIKey" {
		return nil, fmt.Errorf("coordiantes or APIKey was not provided, cannot continue: %v %s", c, lookupMethod)
	}
	// Build for condition with coordiantes.  To be used with HERE API
	if lookupMethod == "Coordinates" {
		monthlyPrayerData.Longitude = c.Coordinates.Longitude
		monthlyPrayerData.Latitude = c.Coordinates.Latitude
		return AladhanData(monthlyPrayerData)
	}
	// Build for condition without coordiantes.  To be used with HERE API
	if lookupMethod == "APIKey" {
		hereLookup.PostalCode = c.PostalCode
		hereLookup.HEREAPIKey = c.HEREAPIKey
		hereResp, hereCity, err := HERECustomerLocation(hereLookup)
		if err != nil {
			return nil, fmt.Errorf("unable to lookup customer location using API Key: %v", hereLookup)
		}

		if hereCity.Coordiantes.Lat == 0 && hereCity.Coordiantes.Lng == 0 {
			for i := 0; i < len(hereResp.Items); i++ {
				if hereResp.Items[i].Address.PostalCode == c.PostalCode {
					monthlyPrayerData.Longitude = hereResp.Items[i].Position.Lng
					monthlyPrayerData.Latitude = hereResp.Items[i].Position.Lat
					return AladhanData(monthlyPrayerData)
				}
			}
			return nil, fmt.Errorf("unable to pinpoint customer location based on zip code: %v:", hereResp)
		}
		monthlyPrayerData.Longitude = hereCity.Coordiantes.Lng
		monthlyPrayerData.Latitude = hereCity.Coordiantes.Lat
		return AladhanData(monthlyPrayerData)
	}

	return nil, fmt.Errorf("unable to locate customer input.  Perhaps not enough input data was given %v:", c)
}

func (c *CustomerLocationInput)checkCustomerInput() (string, error) {
	// Check if API key and Coordinates are not filled
	if c.HEREAPIKey == "" && (c.Coordinates.Longitude == 0 || c.Coordinates.Latitude == 0) {
		_, err := fmt.Fprintf(os.Stderr, "error: HEREAPIKey and coordinates are not filled.  Must fill one or the other")
		return "", err
	}

	// Check if API key and Coordinates are filled
	if c.HEREAPIKey != "" && (c.Coordinates.Longitude != 0 || c.Coordinates.Latitude != 0) {
		_, err := fmt.Fprintf(os.Stderr, "error: HEREAPIKey and Coordinates are filled.  Cannot fill both fields")
		return "", err
	}

	if c.HEREAPIKey != "" && (c.Coordinates.Longitude == 0 || c.Coordinates.Latitude == 0) {
		return "APIKey", nil
	}

	if c.HEREAPIKey == "" && (c.Coordinates.Longitude != 0 || c.Coordinates.Latitude != 0) {
		return "Coordinates", nil
	}
	return "", fmt.Errorf("Could not determine which method to use between API key or Coordinates")
}
