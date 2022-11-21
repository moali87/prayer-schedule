package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

// determineSelectedPrayer tests whether prayerToTest is within current time or after
func DetermineSelectedPrayer(clientTimeNow time.Time, prayerToTest string) (bool, error) {
    fmt.Println(clientTimeNow)
    fmt.Println(prayerToTest)
	prayerToTestHourStr := strings.Split(prayerToTest, ":")[0]
	prayerToTestMinutePre := strings.Split(prayerToTest, ":")[1]
	prayerToTestMinuteStr := strings.Split(prayerToTestMinutePre, "(")[0]
    prayerToTestMinuteStr = strings.TrimSpace(prayerToTestMinuteStr)

	// Convert time data to int where possible
	prayerToTestHour, err := strconv.Atoi(prayerToTestHourStr)
    fmt.Printf("Prayer to test hour Int: %d \n", prayerToTestHour)
	if err != nil {
        return false, fmt.Errorf("unable to convert hour from string to int: %s", err)
	}

	prayerToTestMinute, err := strconv.Atoi(prayerToTestMinuteStr)
    fmt.Printf("Prayer to test minute Int: %d \n", prayerToTestMinute)
	if err != nil {
        return false, fmt.Errorf("unable to convert minute from string to int: %s", err)
	}

	if clientTimeNow.Hour() > prayerToTestHour {
		return true, nil
	}

	if clientTimeNow.Hour() == prayerToTestHour &&
    clientTimeNow.Minute() >= prayerToTestMinute {
        fmt.Printf("Prayer hour and client hour matches. \n")
        fmt.Printf("Client minute is greater than or equal to prayer minute.")
		return true, nil
	}

	if clientTimeNow.Hour() == prayerToTestHour &&
    clientTimeNow.Minute() == prayerToTestMinute {
        fmt.Printf("Prayer hour and client hour matches. \n")
        fmt.Printf("Prayer minute and client minute matches")
		return true, nil
	}

	return false, nil
}

// DeterminedPrayerOutput is the structure which contains current and next prayer name and their time difference
type DeterminedPrayerOutput struct {
	CurrentPrayerName string
	NextPrayerName    string
	PreviousDayIsha   bool
	CurrentPrayerTime string
	NextPrayerTime    string
	TimeDiff          time.Duration
}

// DetermineWhichPrayer returns the current and next prayer.  It will also state if the current prayer is at the previous day
func DetermineWhichPrayer(
	previousDayPrayers *FiveDailyPrayers,
	currentDayPrayers *FiveDailyPrayers,
	nextDayPrayers *FiveDailyPrayers,
	clientTimeNow *time.Time) (*DeterminedPrayerOutput, error) {

	output := &DeterminedPrayerOutput{}
	// Determine if current prayer is Isha for the current day
    determineIsha, err := DetermineSelectedPrayer(*clientTimeNow, currentDayPrayers.Isha)
    if err != nil {
        return nil, err
    }
	if determineIsha {
		output.CurrentPrayerName = "Isha"
		output.NextPrayerName = "Fajr"
		output.PreviousDayIsha = false
		return output, nil
	}

	var currentDayPrayersMap = make(map[string]string)
	var previousDayPrayersMap = make(map[string]string)
	var nextDayPrayersMap = make(map[string]string)
	currentPrayerConvErr := mapstructure.Decode(currentDayPrayers, &currentDayPrayersMap)
	if currentPrayerConvErr != nil {
		return nil, fmt.Errorf("unable to convert current day prayers into a map: %s", currentPrayerConvErr.Error())
	}
	previousPrayerConvErr := mapstructure.Decode(previousDayPrayers, &previousDayPrayersMap)
	if previousPrayerConvErr != nil {
		return nil, fmt.Errorf("unable to convert previous day prayers into a map: %s", previousPrayerConvErr.Error())
	}

	nextPrayerConvErr := mapstructure.Decode(nextDayPrayers, &nextDayPrayersMap)
	if nextPrayerConvErr != nil {
		return nil, fmt.Errorf("unable to convert next day prayers into a map: %s", nextPrayerConvErr.Error())
	}

	// Test if all selected prayers are false.  In which case it's the previous day Isha
	prayerMapKeys := make([]string, len(currentDayPrayersMap))
	prayerMapKeys[0] = "Fajr"
	prayerMapKeys[1] = "Sunrise"
	prayerMapKeys[2] = "Dhuhr"
	prayerMapKeys[3] = "Asr"
	prayerMapKeys[4] = "Maghrib"
	prayerMapKeys[5] = "Isha"

	for key := 0; key < len(prayerMapKeys); key++ {
		nextPrayerName := prayerMapKeys[key]
		selectedPrayerDetermined, err := DetermineSelectedPrayer(*clientTimeNow, currentDayPrayersMap[nextPrayerName])
		if nextPrayerName == "Fajr" && !selectedPrayerDetermined {
			output.CurrentPrayerName = "Isha"
			output.CurrentPrayerTime = currentDayPrayersMap[output.CurrentPrayerName]
			output.PreviousDayIsha = true
			output.NextPrayerName = "Fajr"
			output.NextPrayerTime = nextDayPrayersMap[output.NextPrayerName]
			break
		}
        if err != nil {
            return nil, err
        }

		if nextPrayerName == "Sunrise" && !selectedPrayerDetermined {
			output.CurrentPrayerName = "Fajr"
			output.CurrentPrayerTime = currentDayPrayersMap[output.CurrentPrayerName]
			output.NextPrayerName = "Dhuhr"
			output.PreviousDayIsha = false
			output.NextPrayerTime = nextDayPrayersMap[output.NextPrayerName]
			break
		}

		if nextPrayerName == "Dhuhr" && !selectedPrayerDetermined {
			output.CurrentPrayerName = "Sunrise"
			output.CurrentPrayerTime = currentDayPrayersMap[output.CurrentPrayerName]
			output.NextPrayerName = nextPrayerName
			output.PreviousDayIsha = false
			output.NextPrayerTime = nextDayPrayersMap[output.NextPrayerName]
			break
		}

		if nextPrayerName == "Asr" && !selectedPrayerDetermined {
			output.CurrentPrayerName = "Dhuhr"
			output.CurrentPrayerTime = currentDayPrayersMap[output.CurrentPrayerName]
			output.NextPrayerName = nextPrayerName
			output.PreviousDayIsha = false
			output.NextPrayerTime = nextDayPrayersMap[output.NextPrayerName]
			break
		}

		if nextPrayerName == "Maghrib" && !selectedPrayerDetermined {
			output.CurrentPrayerName = "Asr"
			output.CurrentPrayerTime = currentDayPrayersMap[output.CurrentPrayerName]
			output.NextPrayerName = nextPrayerName
			output.PreviousDayIsha = false
			output.NextPrayerTime = nextDayPrayersMap[output.NextPrayerName]
			break
		}

		if nextPrayerName == "Isha" && !selectedPrayerDetermined {
			output.CurrentPrayerName = "Maghrib"
			output.CurrentPrayerTime = currentDayPrayersMap[output.CurrentPrayerName]
			output.NextPrayerName = nextPrayerName
			output.PreviousDayIsha = false
			output.NextPrayerTime = nextDayPrayersMap[output.NextPrayerName]
			break
		}

		// Break on the first false.  If prayer is Asr and Asr selected prayer returns false.  Set current prayer name to Asr and next prayer to Asr + 1
		if !selectedPrayerDetermined {
			output.CurrentPrayerName = nextPrayerName
			output.CurrentPrayerTime = currentDayPrayersMap[output.CurrentPrayerName]
			output.NextPrayerName = prayerMapKeys[key+1]
			output.NextPrayerTime = nextDayPrayersMap[output.NextPrayerName]
			break
		}
	}

	nextPrayerSplit := strings.Split(output.NextPrayerTime, ":")
	nextPrayerHour := nextPrayerSplit[0]
	nextPrayerMinute := nextPrayerSplit[1]
    nextPrayerMinute = strings.Split(nextPrayerMinute, "(")[0]
    nextPrayerMinute = strings.TrimSpace(nextPrayerMinute)
	timediff, err := timeDiff(clientTimeNow, nextPrayerHour, nextPrayerMinute)
	if err != nil {
        return nil, fmt.Errorf("unable to gather time difference: %s", err)
	}
	output.TimeDiff = *timediff

	fmt.Println(output)
	if output != nil {
		return output, nil
	}

	return nil, fmt.Errorf("unable to pinpoint prayer time for next prayer")
}

// timeDiff calculates the time difference between client current time and next prayer time
func timeDiff(clientTimeNow *time.Time, nextPrayerHour string, nextPrayerMinute string) (*time.Duration, error) {
	intNextPrayerHour, err := strconv.Atoi(nextPrayerHour)
	if err != nil {
		return nil, err
	}
	intNextPrayerMinute, err := strconv.Atoi(nextPrayerMinute)
	if err != nil {
		return nil, err
	}

	var clientTimeNowDay int = clientTimeNow.Day()

	if intNextPrayerHour < clientTimeNow.Hour() {
		clientTimeNowDay = clientTimeNow.AddDate(int(clientTimeNow.Year()), int(clientTimeNow.Month()), int(clientTimeNow.Day())+1).Day()
	}

	nextPrayerTime := time.Date(clientTimeNow.Year(), clientTimeNow.Month(), clientTimeNowDay, intNextPrayerHour, intNextPrayerMinute, 0, 0, clientTimeNow.Location())

	diff := nextPrayerTime.Sub(*clientTimeNow)

	return &diff, nil
}
