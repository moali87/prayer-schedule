package schedule_test

import (
	"testing"
	"time"
    psched "github.com/moali87/prayer-schedule"
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
	t1DeterminedTime := psched.DetermineSelectedPrayer(t1CurrentTime, "13:04")
	if t1DeterminedTime {
		t.Error("determined time with same hour but off before minute returned an incorrect result")
	}

	/*
	   Test with same hour but different minute
	   current minute should be after prayer time minute, meaning the prayer has started
	*/
	t2CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 13, 5, 0, 0, timeLocation)
	t2DeterminedTime := psched.DetermineSelectedPrayer(t2CurrentTime, "13:04")
	if !t2DeterminedTime {
		t.Error("determined time with same hour but off after minute returned an incorrect result")
	}

	/*
	   Test with same hour and minute
	   prayer has started
	*/
	t3CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 13, 5, 0, 0, timeLocation)
	t3DeterminedTime := psched.DetermineSelectedPrayer(t3CurrentTime, "13:04")
	if !t3DeterminedTime {
		t.Error("determined time with same hour and same minute returned an incorrect result")
	}

	/*
	   Test with same different hour and same minute
	   Hour should be before prayer time hour, meaning the prayer has not started yet
	*/
	t4CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 12, 5, 0, 0, timeLocation)
	t4DeterminedTime := psched.DetermineSelectedPrayer(t4CurrentTime, "13:04")
	if t4DeterminedTime {
		t.Error("determined time with different before hour and same minute returned an incorrect result")
	}

	/*
	   Test with same different hour and same minute
	   Hour should be after prayer time hour, meaning the prayer has started
	*/
	t5CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 14, 5, 0, 0, timeLocation)
	t5DeterminedTime := psched.DetermineSelectedPrayer(t5CurrentTime, "13:04")
	if !t5DeterminedTime {
		t.Error("determined time with different after hour and same minute returned an incorrect result")
	}

	/*
	   Test both hour and minute is before prayer hour and minute
	*/
	t6CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 12, 0, 0, 0, timeLocation)
	t6DeterminedTime := psched.DetermineSelectedPrayer(t6CurrentTime, "13:04")
	if t6DeterminedTime {
		t.Error("determined time with different after hour and same minute returned an incorrect result")
	}

	/*
	   Test both hour and minute is after prayer hour and minute
	*/
	t7CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 10, 0, 0, timeLocation)
	t7DeterminedTime := psched.DetermineSelectedPrayer(t7CurrentTime, "13:04")
	if !t7DeterminedTime {
		t.Error("determined time with different after hour and same minute returned an incorrect result")
	}
}

var prevDayPrayerStruct = &psched.FiveDailyPrayers{
	Fajr:    "04:37",
	Sunrise: "06:36",
	Dhuhr:   "13:04",
	Asr:     "16:37",
	Maghrib: "19:34",
	Isha:    "21:32",
}

var currDayPrayerStruct = &psched.FiveDailyPrayers{
	Fajr:    "04:37",
	Sunrise: "06:36",
	Dhuhr:   "13:04",
	Asr:     "16:37",
	Maghrib: "19:34",
	Isha:    "21:33",
}

var nextDayPrayerStruct = &psched.FiveDailyPrayers{
	Fajr:    "04:38",
	Sunrise: "06:36",
	Dhuhr:   "13:04",
	Asr:     "16:37",
	Maghrib: "19:34",
	Isha:    "21:33",
}

// TODO: Create test for determineWhichPrayer function
func TestDetermineWhichPrayerIsha(t *testing.T) {
	timeLocation, err := time.LoadLocation("Local")
	if err != nil {
		t.Errorf("unable to load time location: %s", err.Error())
	}

	// Test if current prayer is previous day Isha
	t1CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 4, 36, 0, 0, timeLocation)
	prevDayIshaStruct, err := psched.DetermineWhichPrayer(prevDayPrayerStruct, currDayPrayerStruct, nextDayPrayerStruct, &t1CurrentTime)
	if err != nil {
		t.Errorf("unable to determine current, previous, and next prayer time structure: %s", err.Error())
	}

	//// Test if current prayer name is Isha
	if prevDayIshaStruct.CurrentPrayerName != "Isha" {
		t.Errorf("incorrect prayer name output.  Prayer name is not \"Isha\", current prayer name is %s", prevDayIshaStruct.CurrentPrayerName)
	}
	//// Test if previous day Isha is true
	if !prevDayIshaStruct.PreviousDayIsha {
		t.Error("previous day Isha is not true when it should be")
	}
	//// Test if next prayer name is Fajr
	if prevDayIshaStruct.NextPrayerName != "Fajr" {
		t.Errorf("incorrect next prayer name. Prayer name is not \"Fajr\", next prayer name is %s", prevDayIshaStruct.NextPrayerName)
	}

	// Test if current prayer is current day Isha
	t2CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 21, 36, 0, 0, timeLocation)
	currDayIshaStruct, err := psched.DetermineWhichPrayer(prevDayPrayerStruct, currDayPrayerStruct, nextDayPrayerStruct, &t2CurrentTime)

	//// Test if current prayer name is Isha
	if currDayIshaStruct.CurrentPrayerName != "Isha" {
		t.Errorf("incorrect prayer name output.  Prayer name is not \"Isha\", current prayer name is %s", currDayIshaStruct.CurrentPrayerName)
	}
	//// Test if currious day Isha is false
	if currDayIshaStruct.PreviousDayIsha {
		t.Error("currious day Isha is not true when it should be")
	}
	//// Test if next prayer name is Fajr
	if currDayIshaStruct.NextPrayerName != "Fajr" {
		t.Errorf("incorrect next prayer name. Prayer name is not \"Fajr\", next prayer name is %s", currDayIshaStruct.NextPrayerName)
	}
}

// Test if currnet prayer name is Fajr and next prayer name is Dhuhr and not sunrise
func TestDetermineWhichPrayerFajr(t *testing.T) {
	timeLocation, err := time.LoadLocation("Local")
	if err != nil {
		t.Errorf("unable to load time location: %s", err.Error())
	}

	t1CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 4, 37, 0, 0, timeLocation)
	currDayFajrStruct, err := psched.DetermineWhichPrayer(prevDayPrayerStruct, currDayPrayerStruct, nextDayPrayerStruct, &t1CurrentTime)
	if err != nil {
		t.Errorf("unable to determine current, previous, and next prayer time structure: %s", err.Error())
	}

	// Test if current prayer name is Fajr
	if currDayFajrStruct.CurrentPrayerName != "Fajr" {
		t.Errorf("incorrect current prayer name. Current day fajr test did not return Fajr as current prayer name %s", currDayFajrStruct.CurrentPrayerName)
	}
	// Test if next prayer name is Dhuhr
	if currDayFajrStruct.NextPrayerName != "Dhuhr" {
		t.Errorf("incorrect next prayer name. Current day fajr test did not return Dhuhr as next prayer name %s", currDayFajrStruct.NextPrayerName)
	}
}

// Test if current prayer name is Sunrise and next prayer name is Dhuhr
func TestDetermineWhichPrayerSunrise(t *testing.T) {
	timeLocation, err := time.LoadLocation("Local")
	if err != nil {
		t.Errorf("unable to load time location: %s", err.Error())
	}

	t1CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 6, 37, 0, 0, timeLocation)
	currDayPrayerStruct, err := psched.DetermineWhichPrayer(prevDayPrayerStruct, currDayPrayerStruct, nextDayPrayerStruct, &t1CurrentTime)
	if err != nil {
		t.Errorf("unable to determine current, previous, and next prayer time structure: %s", err.Error())
	}

	// Test if current prayer name is Fajr
	if currDayPrayerStruct.CurrentPrayerName != "Sunrise" {
		t.Errorf("incorrect current prayer name. Current day sunrise test did not return Sunrise as current prayer name %s", currDayPrayerStruct.CurrentPrayerName)
	}
	// Test if next prayer name is Dhuhr
	if currDayPrayerStruct.NextPrayerName != "Dhuhr" {
		t.Errorf("incorrect next prayer name. Current day sunrise test did not return Dhuhr as next prayer name %s", currDayPrayerStruct.NextPrayerName)
	}
}

// Test if current prayer name is Dhuhr and next prayer name is Asr
func TestDetermineWhichPrayerDhuhr(t *testing.T) {
	timeLocation, err := time.LoadLocation("Local")
	if err != nil {
		t.Errorf("unable to load time location: %s", err.Error())
	}

	t1CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 13, 37, 0, 0, timeLocation)
	currDayDhuhrStruct, err := psched.DetermineWhichPrayer(prevDayPrayerStruct, currDayPrayerStruct, nextDayPrayerStruct, &t1CurrentTime)
	if err != nil {
		t.Errorf("unable to determine current, previous, and next prayer time structure: %s", err.Error())
	}

	// Test if current prayer name is Fajr
	if currDayDhuhrStruct.CurrentPrayerName != "Dhuhr" {
		t.Errorf("incorrect current prayer name. Current day dhuhr test did not return Dhuhr as current prayer name %s", currDayDhuhrStruct.CurrentPrayerName)
	}
	// Test if next prayer name is Dhuhr
	if currDayDhuhrStruct.NextPrayerName != "Asr" {
		t.Errorf("incorrect next prayer name. Current day dhuhr test did not return Asr as next prayer name %s", currDayDhuhrStruct.NextPrayerName)
	}
}

// Test if current prayer name is Asr and next prayer name is Maghrib
func TestDetermineWhichPrayerAsr(t *testing.T) {
	timeLocation, err := time.LoadLocation("Local")
	if err != nil {
		t.Errorf("unable to load time location: %s", err.Error())
	}

	t1CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 16, 37, 0, 0, timeLocation)
	currDayPrayerStruct, err := psched.DetermineWhichPrayer(prevDayPrayerStruct, currDayPrayerStruct, nextDayPrayerStruct, &t1CurrentTime)
	if err != nil {
		t.Errorf("unable to determine current, previous, and next prayer time structure: %s", err.Error())
	}

	// Test if current prayer name is Fajr
	if currDayPrayerStruct.CurrentPrayerName != "Asr" {
		t.Errorf("incorrect current prayer name. Current day asr test did not return Asr as current prayer name %s", currDayPrayerStruct.CurrentPrayerName)
	}
	// Test if next prayer name is Dhuhr
	if currDayPrayerStruct.NextPrayerName != "Maghrib" {
		t.Errorf("incorrect next prayer name. Current day asr test did not return Maghrib as next prayer name %s", currDayPrayerStruct.NextPrayerName)
	}
}

// Test if current prayer name is Asr and next prayer name is Maghrib
func TestDetermineWhichPrayerMaghrib(t *testing.T) {
	timeLocation, err := time.LoadLocation("Local")
	if err != nil {
		t.Errorf("unable to load time location: %s", err.Error())
	}

	t1CurrentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 19, 37, 0, 0, timeLocation)
	currDayPrayerStruct, err := psched.DetermineWhichPrayer(prevDayPrayerStruct, currDayPrayerStruct, nextDayPrayerStruct, &t1CurrentTime)
	if err != nil {
		t.Errorf("unable to determine current, previous, and next prayer time structure: %s", err.Error())
	}

	// Test if current prayer name is Fajr
	if currDayPrayerStruct.CurrentPrayerName != "Maghrib" {
		t.Errorf("incorrect current prayer name. Current day asr test did not return Maghrib as current prayer name %s", currDayPrayerStruct.CurrentPrayerName)
	}
	// Test if next prayer name is Dhuhr
	if currDayPrayerStruct.NextPrayerName != "Isha" {
		t.Errorf("incorrect next prayer name. Current day asr test did not return Isha as next prayer name %s", currDayPrayerStruct.NextPrayerName)
	}
}
