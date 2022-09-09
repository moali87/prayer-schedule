// Package schedule returns customer monthly prayer schedule
package schedule

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// CustomerLocationInputWithHEREAPIKey is a struct which contains data to lookup customer data to the nearest city
type CustomerLocationInputWithHEREAPIKey struct {
	HEREAPIKey  string
	City        string
	CountryCode string
	PostalCode  string
}

// CustomerCoordinatesOutput contains customer coordinates to the nearest city
type CustomerCoordinatesOutput struct {
	Lng float32 `json:"Longitude"`
	Lat float32 `json:"Latitude"`
}

// HERECustomerCityAddressOutput is the output containing total customer location to the nearest city
type HERECustomerCityAddressOutput struct {
	Country     string
	City        string
	PostalCode  string
	Coordiantes CustomerCoordinatesOutput
}

// HERECustomerLocationOutput is the general output which is decoded by JSON from HERE url request
type HERECustomerLocationOutput struct {
	StatusCode int
	Response   struct {
		View []struct {
			Result []struct {
				Location struct {
					DisplayPosition CustomerCoordinatesOutput
					Address         HERECustomerCityAddressOutput
				}
			}
		}
	}
}

// HERECustomerLocation Returns customer location data to the nearest city
func HERECustomerLocation(hereRequestParamaters *CustomerLocationInputWithHEREAPIKey) (*HERECustomerLocationOutput, *HERECustomerCityAddressOutput, error) {
	const hereRestAPI = "https://geocoder.ls.hereapi.com/6.2/geocode.json"
	reqURL := fmt.Sprintf(
		"%s?apiKey=%s&city=%s&countryCode=%s&postalCode=%s",
		hereRestAPI,
		hereRequestParamaters.HEREAPIKey,
		url.QueryEscape(hereRequestParamaters.City),
		hereRequestParamaters.CountryCode,
		hereRequestParamaters.PostalCode,
	)

	resp := new(HERECustomerLocationOutput)
	req, err := http.Get(reqURL)
	if err != nil {
		errMsg := fmt.Errorf("Unable to retrieve customer location")
		return nil, nil, errMsg
	}

	defer func() {
		err := req.Body.Close()
		if err != nil {
			fmt.Println("Failed to close request from HERE API")
			panic(err)
		}
	}()

	json.NewDecoder(req.Body).Decode(resp)
	resp.StatusCode = req.StatusCode
	if req.StatusCode != 200 {
		fmt.Printf("HERE API response is not 200: %v", req.StatusCode)
		fmt.Println(resp)
		return nil, nil, fmt.Errorf("Return code not 200: %d", req.StatusCode)
	}

	HERECustomerCityAddressOutputStruct := new(HERECustomerCityAddressOutput)

	for i := 0; i < len(resp.Response.View[0].Result); i++ {
		if resp.Response.View[0].Result[i].Location.Address.PostalCode == hereRequestParamaters.PostalCode {
			HERECustomerAddress := resp.Response.View[0].Result[i].Location
			HERECustomerCityAddressOutputStruct.City = HERECustomerAddress.Address.City
			HERECustomerCityAddressOutputStruct.Country = HERECustomerAddress.Address.Country
			HERECustomerCityAddressOutputStruct.PostalCode = HERECustomerAddress.Address.PostalCode
			HERECustomerCityAddressOutputStruct.Coordiantes.Lat = HERECustomerAddress.DisplayPosition.Lat
			HERECustomerCityAddressOutputStruct.Coordiantes.Lng = HERECustomerAddress.DisplayPosition.Lng
		}
	}

	return resp, HERECustomerCityAddressOutputStruct, nil
}
