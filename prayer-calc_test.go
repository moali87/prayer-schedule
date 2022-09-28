package schedule

import (
	"testing"
	"time"
)

func TestDetermineSelectedPrayer(t *testing.T) {
  timeLocation, err := time.LoadLocation("Local")
  if err != nil {
    t.Errorf("unable to load time location: %s", err.Error())
  }

  /* 
    Test with same hour but different minute
    current minute should be before prayer time minute, meaning the prayer has not started yet
  */
  t1CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 13, 0, 0, 0, timeLocation)
  t1DeterminedTime := determineSelectedPrayer(t1CurrentTime, "13:04")  
  if t1DeterminedTime {
    t.Error("determined time with same hour but off before minute returned an incorrect result")
  }

  /* 
    Test with same hour but different minute
    current minute should be after prayer time minute, meaning the prayer has started
  */
  t2CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 13, 5, 0, 0, timeLocation)
  t2DeterminedTime := determineSelectedPrayer(t2CurrentTime, "13:04")  
  if !t2DeterminedTime {
    t.Error("determined time with same hour but off after minute returned an incorrect result")
  }
  
  /* 
    Test with same hour and minute
    prayer has started
  */
  t3CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 13, 5, 0, 0, timeLocation)
  t3DeterminedTime := determineSelectedPrayer(t3CurrentTime, "13:04")  
  if !t3DeterminedTime {
    t.Error("determined time with same hour and same minute returned an incorrect result")
  }

  /* 
    Test with same different hour and same minute
    Hour should be before prayer time hour, meaning the prayer has not started yet
  */
  t4CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 12, 5, 0, 0, timeLocation)
  t4DeterminedTime := determineSelectedPrayer(t4CurrentTime, "13:04")  
  if t4DeterminedTime {
    t.Error("determined time with different before hour and same minute returned an incorrect result")
  }

  /* 
    Test with same different hour and same minute
    Hour should be after prayer time hour, meaning the prayer has started
  */
  t5CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 14, 5, 0, 0, timeLocation)
  t5DeterminedTime := determineSelectedPrayer(t5CurrentTime, "13:04")  
  if !t5DeterminedTime {
    t.Error("determined time with different after hour and same minute returned an incorrect result")
  }

  /*
    Test both hour and minute is before prayer hour and minute
  */
  t6CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 12, 0, 0, 0, timeLocation)
  t6DeterminedTime := determineSelectedPrayer(t6CurrentTime, "13:04")  
  if t6DeterminedTime {
    t.Error("determined time with different after hour and same minute returned an incorrect result")
  }

  /*
    Test both hour and minute is after prayer hour and minute
  */
  t7CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 10, 0, 0, timeLocation)
  t7DeterminedTime := determineSelectedPrayer(t7CurrentTime, "13:04")  
  if !t7DeterminedTime {
    t.Error("determined time with different after hour and same minute returned an incorrect result")
  }
} 

var prevDayPrayerStruct = &FiveDailyPrayers{
  Fajr: "04:37",
  Sunrise: "06:36",
  Dhuhr: "13:04",
  Asr: "16:37",
  Maghrib: "19:34",
  Isha: "21:32",
}

var currDayPrayerStruct = &FiveDailyPrayers{
  Fajr: "04:37",
  Sunrise: "06:36",
  Dhuhr: "13:04",
  Asr: "16:37",
  Maghrib: "19:34",
  Isha: "21:33",
}

var nextDayPrayerStruct = &FiveDailyPrayers{
  Fajr: "04:38",
  Sunrise: "06:36",
  Dhuhr: "13:04",
  Asr: "16:37",
  Maghrib: "19:34",
  Isha: "21:33",
}

// TODO: Create test for determineWhichPrayer function
func TestDetermineWhichPrayerIsha(t *testing.T) {
  timeLocation, err := time.LoadLocation("Local")
  if err != nil {
    t.Errorf("unable to load time location: %s", err.Error())
  }

  // Test if current prayer is previous day Isha 
  t1CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 4, 36, 0, 0, timeLocation)
  prevDayIshaStruct, err := determineWhichPrayer(prevDayPrayerStruct, currDayPrayerStruct, nextDayPrayerStruct, &t1CurrentTime)
  if err != nil {
    t.Errorf("unable to determine current, previous, and next prayer time structure: %s", err.Error())
  }

  //// Test if current prayer name is Isha
  if prevDayIshaStruct.currentPrayerName != "Isha" {
    t.Errorf("incorrect prayer name output.  Prayer name is not \"Isha\", current prayer name is %s", prevDayIshaStruct.currentPrayerName)
  }
  //// Test if previous day Isha is true
  if !prevDayIshaStruct.previousDayIsha {
    t.Error("previous day Isha is not true when it should be")
  }
  //// Test if next prayer name is Fajr
  if prevDayIshaStruct.nextPrayerName != "Fajr" {
    t.Errorf("incorrect next prayer name. Prayer name is not \"Fajr\", next prayer name is %s", prevDayIshaStruct.nextPrayerName)
  }

  // Test if current prayer is current day Isha 
  t2CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 21, 36, 0, 0, timeLocation)
  currDayIshaStruct, err := determineWhichPrayer(prevDayPrayerStruct, currDayPrayerStruct, nextDayPrayerStruct, &t2CurrentTime)

  //// Test if current prayer name is Isha
  if currDayIshaStruct.currentPrayerName != "Isha" {
    t.Errorf("incorrect prayer name output.  Prayer name is not \"Isha\", current prayer name is %s", currDayIshaStruct.currentPrayerName)
  }
  //// Test if currious day Isha is false
  if currDayIshaStruct.previousDayIsha {
    t.Error("currious day Isha is not true when it should be")
  }
  //// Test if next prayer name is Fajr
  if currDayIshaStruct.nextPrayerName != "Fajr" {
    t.Errorf("incorrect next prayer name. Prayer name is not \"Fajr\", next prayer name is %s", currDayIshaStruct.nextPrayerName)
  }
}

// Test if currnet prayer name is Fajr and next prayer name is Dhuhr and not sunrise
func TestDetermineWhichPrayerFajr(t *testing.T) {
  timeLocation, err := time.LoadLocation("Local")
  if err != nil {
    t.Errorf("unable to load time location: %s", err.Error())
  }

  t1CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 4, 37, 0, 0, timeLocation)
  currDayFajrStruct, err := determineWhichPrayer(prevDayPrayerStruct, currDayPrayerStruct, nextDayPrayerStruct, &t1CurrentTime)
  if err != nil {
    t.Errorf("unable to determine current, previous, and next prayer time structure: %s", err.Error())
  }

  // Test if current prayer name is Fajr
  if currDayFajrStruct.currentPrayerName != "Fajr" {
    t.Errorf("incorrect current prayer name. Current day fajr test did not return Fajr as current prayer name %s", currDayFajrStruct.currentPrayerName)
  }
  // Test if next prayer name is Dhuhr
  if currDayFajrStruct.nextPrayerName != "Dhuhr" {
    t.Errorf("incorrect next prayer name. Current day fajr test did not return Dhuhr as next prayer name %s", currDayFajrStruct.nextPrayerName)
  }
}

// Test if current prayer name is Sunrise and next prayer name is Dhuhr
func TestDetermineWhichPrayerSunrise(t *testing.T) {
  timeLocation, err := time.LoadLocation("Local")
  if err != nil {
    t.Errorf("unable to load time location: %s", err.Error())
  }

  t1CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 6, 37, 0, 0, timeLocation)
  currDayPrayerStruct, err := determineWhichPrayer(prevDayPrayerStruct, currDayPrayerStruct, nextDayPrayerStruct, &t1CurrentTime)
  if err != nil {
    t.Errorf("unable to determine current, previous, and next prayer time structure: %s", err.Error())
  }

  // Test if current prayer name is Fajr
  if currDayPrayerStruct.currentPrayerName != "Sunrise" {
    t.Errorf("incorrect current prayer name. Current day sunrise test did not return Sunrise as current prayer name %s", currDayPrayerStruct.currentPrayerName)
  }
  // Test if next prayer name is Dhuhr
  if currDayPrayerStruct.nextPrayerName != "Dhuhr" {
    t.Errorf("incorrect next prayer name. Current day sunrise test did not return Dhuhr as next prayer name %s", currDayPrayerStruct.nextPrayerName)
  }
}

// Test if current prayer name is Dhuhr and next prayer name is Asr
func TestDetermineWhichPrayerDhuhr(t *testing.T) {
  timeLocation, err := time.LoadLocation("Local")
  if err != nil {
    t.Errorf("unable to load time location: %s", err.Error())
  }

  t1CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 13, 37, 0, 0, timeLocation)
  currDayDhuhrStruct, err := determineWhichPrayer(prevDayPrayerStruct, currDayPrayerStruct, nextDayPrayerStruct, &t1CurrentTime)
  if err != nil {
    t.Errorf("unable to determine current, previous, and next prayer time structure: %s", err.Error())
  }

  // Test if current prayer name is Fajr
  if currDayDhuhrStruct.currentPrayerName != "Dhuhr" {
    t.Errorf("incorrect current prayer name. Current day dhuhr test did not return Dhuhr as current prayer name %s", currDayDhuhrStruct.currentPrayerName)
  }
  // Test if next prayer name is Dhuhr
  if currDayDhuhrStruct.nextPrayerName != "Asr" {
    t.Errorf("incorrect next prayer name. Current day dhuhr test did not return Asr as next prayer name %s", currDayDhuhrStruct.nextPrayerName)
  }
}

// Test if current prayer name is Asr and next prayer name is Maghrib
func TestDetermineWhichPrayerAsr(t *testing.T) {
  timeLocation, err := time.LoadLocation("Local")
  if err != nil {
    t.Errorf("unable to load time location: %s", err.Error())
  }

  t1CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 16, 37, 0, 0, timeLocation)
  currDayPrayerStruct, err := determineWhichPrayer(prevDayPrayerStruct, currDayPrayerStruct, nextDayPrayerStruct, &t1CurrentTime)
  if err != nil {
    t.Errorf("unable to determine current, previous, and next prayer time structure: %s", err.Error())
  }

  // Test if current prayer name is Fajr
  if currDayPrayerStruct.currentPrayerName != "Asr" {
    t.Errorf("incorrect current prayer name. Current day asr test did not return Asr as current prayer name %s", currDayPrayerStruct.currentPrayerName)
  }
  // Test if next prayer name is Dhuhr
  if currDayPrayerStruct.nextPrayerName != "Maghrib" {
    t.Errorf("incorrect next prayer name. Current day asr test did not return Maghrib as next prayer name %s", currDayPrayerStruct.nextPrayerName)
  }
}

// Test if current prayer name is Asr and next prayer name is Maghrib
func TestDetermineWhichPrayerMaghrib(t *testing.T) {
  timeLocation, err := time.LoadLocation("Local")
  if err != nil {
    t.Errorf("unable to load time location: %s", err.Error())
  }

  t1CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 19, 37, 0, 0, timeLocation)
  currDayPrayerStruct, err := determineWhichPrayer(prevDayPrayerStruct, currDayPrayerStruct, nextDayPrayerStruct, &t1CurrentTime)
  if err != nil {
    t.Errorf("unable to determine current, previous, and next prayer time structure: %s", err.Error())
  }

  // Test if current prayer name is Fajr
  if currDayPrayerStruct.currentPrayerName != "Maghrib" {
    t.Errorf("incorrect current prayer name. Current day asr test did not return Maghrib as current prayer name %s", currDayPrayerStruct.currentPrayerName)
  }
  // Test if next prayer name is Dhuhr
  if currDayPrayerStruct.nextPrayerName != "Isha" {
    t.Errorf("incorrect next prayer name. Current day asr test did not return Isha as next prayer name %s", currDayPrayerStruct.nextPrayerName)
  }
}
