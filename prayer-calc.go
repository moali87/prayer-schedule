package schedule

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

// determineSelectedPrayer tests whether prayerToTest is within current time or after
func determineSelectedPrayer(clientTimeNow time.Time, prayerToTest string) bool {
  prayerToTestHourStr := strings.Split(prayerToTest, ":")[0]
  prayerToTestMinutePre := strings.Split(prayerToTest, ":")[1]
  prayerToTestMinuteStr := strings.Split(prayerToTestMinutePre, "(")[0]

  // Convert time data to int where possible
  prayerToTestHour, err := strconv.Atoi(prayerToTestHourStr)
  if err != nil {
    log.Fatal("unable to convert minute from string to int")
  }

  prayerToTestMinute, err := strconv.Atoi(prayerToTestMinuteStr)
  if err != nil {
    log.Fatal("unable to convert minute frm string to int")
  }

  if clientTimeNow.Hour() > prayerToTestHour {
    return true
  }

  if clientTimeNow.Hour() == prayerToTestHour && clientTimeNow.Minute() >= prayerToTestMinute {
    return true
  }

  if clientTimeNow.Hour() == prayerToTestHour && clientTimeNow.Minute() == prayerToTestMinute {
    return true
  }
  
  return false
}

type determinedPrayerOutput struct {
  currentPrayerName string
  nextPrayerName string
  previousDayIsha bool
  currentPrayerTime string
  nextPrayerTime string
  timeDiff time.Duration
}

// determineWhichPrayer returns the current and next prayer.  It will also state if the current prayer is at the previous day
func determineWhichPrayer(
  previousDayPrayers *FiveDailyPrayers,
  currentDayPrayers *FiveDailyPrayers,
  nextDayPrayers *FiveDailyPrayers,
  clientTimeNow *time.Time) (*determinedPrayerOutput, error) {

    output := &determinedPrayerOutput{}
    // Determine if current prayer is Isha for the current day
    if determineSelectedPrayer(*clientTimeNow, currentDayPrayers.Isha) {
      output.currentPrayerName = "Isha"
      output.nextPrayerName = "Fajr"
      output.previousDayIsha = false
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
      selectedPrayerDetermined := determineSelectedPrayer(*clientTimeNow, currentDayPrayersMap[nextPrayerName])
      if nextPrayerName == "Fajr" && !selectedPrayerDetermined {
        output.currentPrayerName = "Isha"
        output.currentPrayerTime = currentDayPrayersMap[output.currentPrayerName]
        output.previousDayIsha = true
        output.nextPrayerName = "Fajr"
        output.nextPrayerTime = nextDayPrayersMap[output.nextPrayerName]
        break
      }

      if nextPrayerName == "Sunrise" && !selectedPrayerDetermined {
        output.currentPrayerName = "Fajr"
        output.currentPrayerTime = currentDayPrayersMap[output.currentPrayerName]
        output.nextPrayerName = "Dhuhr"
        output.previousDayIsha = false
        output.nextPrayerTime = nextDayPrayersMap[output.nextPrayerName]
        break
      }

      if nextPrayerName == "Dhuhr" && !selectedPrayerDetermined {
        output.currentPrayerName = "Sunrise"
        output.currentPrayerTime = currentDayPrayersMap[output.currentPrayerName]
        output.nextPrayerName = nextPrayerName
        output.previousDayIsha = false
        output.nextPrayerTime = nextDayPrayersMap[output.nextPrayerName]
        break
      }

      if nextPrayerName == "Asr" && !selectedPrayerDetermined {
        output.currentPrayerName = "Dhuhr"
        output.currentPrayerTime = currentDayPrayersMap[output.currentPrayerName]
        output.nextPrayerName = nextPrayerName
        output.previousDayIsha = false
        output.nextPrayerTime = nextDayPrayersMap[output.nextPrayerName]
        break
      }

      if nextPrayerName == "Maghrib" && !selectedPrayerDetermined {
        output.currentPrayerName = "Asr"
        output.currentPrayerTime = currentDayPrayersMap[output.currentPrayerName]
        output.nextPrayerName = nextPrayerName
        output.previousDayIsha = false
        output.nextPrayerTime = nextDayPrayersMap[output.nextPrayerName]
        break
      }

      if nextPrayerName == "Isha" && !selectedPrayerDetermined {
        output.currentPrayerName = "Maghrib"
        output.currentPrayerTime = currentDayPrayersMap[output.currentPrayerName]
        output.nextPrayerName = nextPrayerName
        output.previousDayIsha = false
        output.nextPrayerTime = nextDayPrayersMap[output.nextPrayerName]
        break
      }

      // Break on the first false.  If prayer is Asr and Asr selected prayer returns false.  Set current prayer name to Asr and next prayer to Asr + 1
      if !selectedPrayerDetermined {
        output.currentPrayerName = nextPrayerName
        output.currentPrayerTime = currentDayPrayersMap[output.currentPrayerName]
        output.nextPrayerName = prayerMapKeys[key + 1]
        output.nextPrayerTime = nextDayPrayersMap[output.nextPrayerName]
        break
      }
    }

    nextPrayerSplit := strings.Split(output.nextPrayerTime, ":")
    nextPrayerHour := nextPrayerSplit[0]
    nextPrayerMinute := nextPrayerSplit[1]
    timediff, err := timeDiff(clientTimeNow, nextPrayerHour, nextPrayerMinute)
    if err != nil {
      return nil, fmt.Errorf("unable to gather time difference")
    }
    output.timeDiff = *timediff

    fmt.Println(output)
    if output != nil {
      return output, nil
    }

   return nil, fmt.Errorf("unable to pinpoint prayer time for next prayer")
  }

  // timeDiff calculates the time difference between client current time and next prayer time
  func timeDiff(clientTimeNow *time.Time, nextPrayerHour string, nextPrayerMinute string) (*time.Duration, error){
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
      clientTimeNowDay = clientTimeNow.AddDate(int(clientTimeNow.Year()), int(clientTimeNow.Month()), int(clientTimeNow.Day()) + 1).Day()
    }

    nextPrayerTime := time.Date(clientTimeNow.Year(), clientTimeNow.Month(), clientTimeNowDay, intNextPrayerHour, intNextPrayerMinute, 0, 0, clientTimeNow.Location())

    diff := nextPrayerTime.Sub(*clientTimeNow)

    return &diff, nil
  }
