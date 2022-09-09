// Package schedule returns cutomer monthly prayer data

package schedule

import (
	"fmt"
	"log"
	"os"
	"time"
)

type CustomerLocationInput struct {
	City        string
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
	city string,
	countryCode string,
	customerTime time.Time,
	postalCode string,
	latitude float32,
	longitude float32) (*CustomerLocationInput, error) {
	return &CustomerLocationInput{
		City:        city,
		CountryCode: countryCode,
		PostalCode:  postalCode,
		Coordinates: PrayerCalendarInputCoordinates{
			Latitude:  latitude,
			Longitude: longitude,
		},
	}, nil
}

func NewPrayerCalendarWithoutCoordiantes(
	city string,
	countryCode string,
	customerTime time.Time,
	postalCode string,
	hereAPIKey string,
) (*CustomerLocationInput, error) {
	return &CustomerLocationInput{
		City:        city,
		CountryCode: countryCode,
		PostalCode:  postalCode,
		HEREAPIKey:  hereAPIKey,
	}, nil
}

/* 
  PrayerCalendar returns customer monthly prayer data with or without customer providing coordinates.
  if customer does not provide coordiantes, they must provide a HERE API Key
*/
func PrayerCalendar(customerInput *CustomerLocationInput) (*PCalOutput, error) {
	lookupMethod, err := checkCustomerInput(customerInput)
	if err != nil {
		log.Fatal("CustomerInputError")
	}

	hereLookup := new(CustomerLocationInputWithHEREAPIKey)
	monthlyPrayerData := new(PCalInput)

	monthlyPrayerData.CustTime = customerInput.CustTime
	monthlyPrayerData.Institution = customerInput.Institution
	hereLookup.City = customerInput.City
	hereLookup.CountryCode = customerInput.CountryCode

  if lookupMethod != "Coordinates" && lookupMethod != "APIKey" {
    log.Fatalf("coordiantes or APIKey was not provided, cannot continue: %v %s", customerInput, lookupMethod)
  }
	// Build for condition with coordiantes.  To be used with HERE API
	if lookupMethod == "Coordinates" {
		monthlyPrayerData.Longitude = customerInput.Coordinates.Longitude
		monthlyPrayerData.Latitude = customerInput.Coordinates.Latitude
    return AladhanData(monthlyPrayerData)
	}
	// Build for condition without coordiantes.  To be used with HERE API
	if lookupMethod == "APIKey" {
		hereLookup.HEREAPIKey = customerInput.HEREAPIKey
    hereResp, hereCity, err := HERECustomerLocation(hereLookup)
    if err != nil {
      log.Fatalf("unable to lookup customer location: %v", hereLookup)
    }

    if hereCity.Coordiantes.Lat == 0 && hereCity.Coordiantes.Lng == 0 {
      hereRespResults := hereResp.Response.View[0].Result
      for i := 0; i < len(hereRespResults); i++ {
        if hereRespResults[i].Location.Address.PostalCode == customerInput.PostalCode {
          monthlyPrayerData.Longitude = hereRespResults[i].Location.DisplayPosition.Lng
          monthlyPrayerData.Latitude = hereRespResults[i].Location.DisplayPosition.Lat
          return AladhanData(monthlyPrayerData)
        }
      }
      log.Fatalf("unable to pinpoint customer location based on zip code: %v:", hereRespResults)
    }
    monthlyPrayerData.Longitude = hereCity.Coordiantes.Lng
    monthlyPrayerData.Latitude = hereCity.Coordiantes.Lat
    return AladhanData(monthlyPrayerData)
	}

  log.Fatalf("unable to locate customer input.  Perhaps not enough input data was given %v:", customerInput)

  return new(PCalOutput), err
}

func checkCustomerInput(customerInput *CustomerLocationInput) (string, error) {
	// Check if API key and Coordinates are not filled
	if customerInput.HEREAPIKey == "" && (customerInput.Coordinates.Longitude == 0 || customerInput.Coordinates.Latitude == 0) {
		_, err := fmt.Fprintf(os.Stderr, "error: HEREAPIKey and coordinates are not filled.  Must fill one or the other")
		return "", err
	}

	// Check if API key and Coordinates are filled
	if customerInput.HEREAPIKey != "" && (customerInput.Coordinates.Longitude != 0 || customerInput.Coordinates.Latitude != 0) {
		_, err := fmt.Fprintf(os.Stderr, "error: HEREAPIKey and Coordinates are filled.  Cannot fill both fields")
		return "", err
	}

	if customerInput.HEREAPIKey != "" && (customerInput.Coordinates.Longitude == 0 || customerInput.Coordinates.Latitude == 0) {
		return "APIKey", nil
	}

	if customerInput.HEREAPIKey == "" && (customerInput.Coordinates.Longitude != 0 || customerInput.Coordinates.Latitude != 0) {
		return "Coordinates", nil
	}
	log.Fatal("Could not determine which method to use between API key or Coordinates")
	return "", fmt.Errorf("Could not determine which method to use between API key or Coordinates")
}
