package schedule

import (
	"os"
	"testing"
)

func TestHEREAPI(t *testing.T) {
	customerLocationInput := &CustomerLocationInputWithHEREAPIKey{
		City:        "Beverly Hills",
		CountryCode: "USA",
		HEREAPIKey:  os.Getenv("HERE_API_KEY"),
		PostalCode:  "90210",
	}

	custLocRet, custCityLocRet, err := HERECustomerLocation(customerLocationInput)
	if err != nil {
		t.Errorf("Customer location lookup failed %s", err.Error())
	}

	if custLocRet.StatusCode != 200 {
		t.Errorf("Customer lookup return code is not 200 %d", custLocRet.StatusCode)
	}

	if custCityLocRet.PostalCode != customerLocationInput.PostalCode {
		t.Error("Error: Failed to match postal code to HERE address output")
	}

	if len(custLocRet.Response.View) < 1 {
		t.Errorf("Customer lookup returned with no locations %v", custLocRet.Response.View)
	}

	if len(custLocRet.Response.View[0].Result) < 1 {
		t.Errorf("Customer lookup result is empty %v", custLocRet.Response.View[0].Result)
	}
}
