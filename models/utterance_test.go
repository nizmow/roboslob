package models

import (
	"testing"
	"time"
)

func TestGetStartAndEndOfDay(t *testing.T) {
	loc, _ := time.LoadLocation("Australia/Sydney")
	timeLocal := time.Date(2010, 5, 5, 13, 5, 11, 0, loc)
	startOfDay, endOfDay := getStartAndEndOfDayInLocation(timeLocal, loc)

	if startOfDay != time.Date(2010, 5, 4, 14, 0, 0, 0, time.UTC) {
		t.Fatalf("startOfDay is %s", startOfDay)
	}

	if endOfDay != time.Date(2010, 5, 5, 13, 59, 59, 59, time.UTC) {
		t.Fatalf("endOfDay is %s", endOfDay)
	}
}
