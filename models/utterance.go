package models

import (
	"log"
	"time"
)

// We kind of try to design things to work with multiple zones,
// but for now let's just pretend everyone's in Sydney.
const timezone string = "Australia/Sydney"

type Utterance struct {
	// db fields
	ID        int64
	UpdatedAt time.Time
	CreatedAt time.Time

	TelegramUserID    int `gorm:"index"`
	ActualUtterance   string
	UtteredAt         time.Time `gorm:"index"`
	UtteredAtTimezone string
}

type DayCount struct {
	Weekday time.Weekday
	Count   int
	Date    time.Time
}

const dayCountSql = `
SELECT
	STRFTIME('%Y%m%d', uttered_at) AS date,
	STRFTIME('%w', uttered_at) AS day_of_week,
	COUNT(*) AS count
FROM utterances
WHERE 
	telegram_user_id = ?
	AND uttered_at BETWEEN ? and ?
ORDER BY uttered_at ASC`

func AddUtterance(actualUtterance string, telegramUserID int) (int64, error) {
	utteranceTime := time.Now().UTC()
	utterance := &Utterance{ActualUtterance: actualUtterance, TelegramUserID: telegramUserID, UtteredAt: utteranceTime, UtteredAtTimezone: timezone}
	result := DB.Create(utterance)
	log.Printf("Adding utterance for %d at %s", telegramUserID, utteranceTime)
	return utterance.ID, result.Error
}

func GetCount(atTime time.Time, telegramUserID int) int64 {
	loc, _ := time.LoadLocation(timezone)
	timeInZone := atTime.In(loc)
	startUTC, endUTC := getStartAndEndOfDayInLocation(timeInZone, loc)
	var count int64
	DB.Model(&Utterance{}).Where(&Utterance{TelegramUserID: telegramUserID}).Where("uttered_at BETWEEN ? and ?", startUTC, endUTC).Count(&count)
	log.Printf("Returning count of %d for %d between %s - %s", count, telegramUserID, startUTC, endUTC)
	return count
}

func GetLastSevenDays(endTime time.Time, telegramUserID int) []DayCount {
	// it's 7 days, so let's just set that here
	const numberOfDays = 7

	type rawDayCount struct {
		Date      string
		DayOfWeek int
		Count     int
	}

	// Dates are annoying, but this is where we want to be for '7 days in the past'
	startTime := endTime.AddDate(0, 0, -(numberOfDays - 1))

	// get the UTC time for the end of the end time, and the start of the start time
	loc, _ := time.LoadLocation(timezone)
	_, endTimeUtc := getStartAndEndOfDayInLocation(endTime, loc)
	startTimeUtc, _ := getStartAndEndOfDayInLocation(startTime, loc)

	var rawDayCountResult []rawDayCount
	DB.Raw(dayCountSql, telegramUserID, startTimeUtc, endTimeUtc).Scan(&rawDayCountResult)

	// We have to do two things here: populate the exported struct, and fill in the
	// gaps in the binning done by SQL. This is probably the least elegant thing
	// possible, but I don't know go enough to do any better right now.
	var dayCount []DayCount
	dateIterator := startTimeUtc
	for i := 0; i < numberOfDays; i++ {
		dateIterator = dateIterator.AddDate(0, 0, 1)
		for _, e := range rawDayCountResult {
			// If we have a match in our source from the db, then use it, otherwise
			// create an empty bin.
			if e.Date == dateIterator.Format("20060102") {
				dayCount = append(dayCount, DayCount{Weekday: time.Weekday(e.DayOfWeek), Count: e.Count, Date: dateIterator})
			} else {
				dayCount = append(dayCount, DayCount{Weekday: dateIterator.Weekday(), Count: 0, Date: dateIterator})
			}
		}
	}

	return dayCount
}

func getStartAndEndOfDayInLocation(utcTime time.Time, loc *time.Location) (time.Time, time.Time) {
	timeInZone := utcTime.In(loc)
	y, m, d := timeInZone.Date()
	startOfDayInZone := time.Date(y, m, d, 0, 0, 0, 0, loc)
	endOfDayInZone := time.Date(y, m, d, 23, 59, 59, 59, loc)

	return startOfDayInZone.UTC(), endOfDayInZone.UTC()
}
