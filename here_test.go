package schedule

import (
	"os"
	"testing"
)

func TestHEREAPI(t *testing.T) {
	customerLocationInput := &CustomerLocationInputWithHEREAPIKey{
		CountryCode: "US",
		HEREAPIKey:  os.Getenv("HERE_API_KEY"),
		PostalCode:  "90210",
	}

	custLocRet, _, err := HERECustomerLocation(customerLocationInput)
	if err != nil {
		t.Errorf("Customer location lookup failed %s", err.Error())
	}

	if custLocRet.StatusCode != 200 {
		t.Errorf("Customer lookup return code is not 200 %d", custLocRet.StatusCode)
	}

	if len(custLocRet.Items) < 1 {
		t.Errorf("Customer lookup returned with no locations %v", custLocRet.Items)
	}
}
