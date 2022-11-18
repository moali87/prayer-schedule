package schedule_test

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	psched "github.com/moali87/prayer-schedule"
)

func TestHEREAPI(t *testing.T) {
	customerLocationInput := &psched.CustomerLocationInputWithHEREAPIKey{
		CountryCode: "US",
		HEREAPIKey:  os.Getenv("HERE_API_KEY"),
		PostalCode:  "90210",
	}

    _, filename, _, ok := runtime.Caller(0)

    if !ok {
		panic("No caller information")
	}
	fmt.Printf("Filename : %q, Dir : %q\n", filename, path.Dir(filename))

	custLocRet, _, err := psched.HERECustomerLocation(customerLocationInput)
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
