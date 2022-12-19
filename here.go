// Package schedule returns customer monthly prayer schedule
package schedule

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// CustomerLocationInputWithHEREAPIKey is a struct which contains data to lookup customer data to the nearest city
type CustomerLocationInputWithHEREAPIKey struct {
	HEREAPIKey  string
	CountryCode string
	PostalCode  string
}

// CustomerCoordinatesOutput contains customer coordinates to the nearest city
type CustomerCoordinatesOutput struct {
	Lng float32 `json:"Lng"`
	Lat float32 `json:"Lat"`
}

// HERECustomerCityAddressOutput is the output containing total customer location to the nearest city
type HERECustomerCityAddressOutput struct {
	Country     string
	PostalCode  string
	Coordiantes CustomerCoordinatesOutput
}

type HERECustomerCityAddressOutputAddressLabel struct {
	Label       string `json:"Label"`
	CountryCode string `json:"countryCode"`
	PostalCode  string `json:"postalCode"`
}

// HERECustomerLocationOutput is the general output which is decoded by JSON from HERE url request
type HERECustomerLocationOutput struct {
	StatusCode int
	Items      []struct {
		Address  HERECustomerCityAddressOutputAddressLabel
		Title    string `json:"title"`
		Position CustomerCoordinatesOutput
	}
}

// HERECustomerLocation Returns customer location data to the nearest city
func HERECustomerLocation(hereRequestParamaters *CustomerLocationInputWithHEREAPIKey) (*HERECustomerLocationOutput, *HERECustomerCityAddressOutput, error) {
	resp := new(HERECustomerLocationOutput)
    var countryCode string
    countryCode = hereRequestParamaters.CountryCode

	const hereRestAPI = "https://geocode.search.hereapi.com/v1/geocode"
	reqURL := fmt.Sprintf(
		"%s?in=countryCode:%s&qq=postalCode=%s&apiKey=%s",
		hereRestAPI,
		strings.ToUpper(countryCode),
		hereRequestParamaters.PostalCode,
		hereRequestParamaters.HEREAPIKey,
	)

	req, err := http.Get(reqURL)
	if err != nil {
		errMsg := fmt.Errorf("Unable to retrieve customer location")
		return resp, nil, errMsg
	}

	defer func() {
		err := req.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	json.NewDecoder(req.Body).Decode(resp)
	resp.StatusCode = req.StatusCode
    fmt.Printf("HERE rest API response code: %d", req.StatusCode)
	if req.StatusCode != 200 {
		fmt.Printf("HERE API response is not 200: %d", req.StatusCode)
		fmt.Println(resp)
		return resp, nil, fmt.Errorf("Return code not 200: %d", req.StatusCode)
	}

	HERECustomerCityAddressOutputStruct := new(HERECustomerCityAddressOutput)

	for i := 0; i < len(resp.Items); i++ {
		if resp.Items[i].Address.PostalCode == hereRequestParamaters.PostalCode {
			HERECustomerAddress := resp.Items[i]
			HERECustomerCityAddressOutputStruct.Country = HERECustomerAddress.Address.CountryCode
			HERECustomerCityAddressOutputStruct.PostalCode = HERECustomerAddress.Address.PostalCode
			HERECustomerCityAddressOutputStruct.Coordiantes.Lat = HERECustomerAddress.Position.Lat
			HERECustomerCityAddressOutputStruct.Coordiantes.Lng = HERECustomerAddress.Position.Lng
		}
	}

    if len(resp.Items) == 0 {
        return nil, nil, fmt.Errorf("HERE did not find any addresses")
    }

	return resp, HERECustomerCityAddressOutputStruct, nil
}
