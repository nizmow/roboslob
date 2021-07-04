package models

import "time"

var timezone string = "Australia/Sydney"

type Utterance struct {
	// db fields
	ID        int64
	UpdatedAt time.Time
	CreatedAt time.Time

	TelegramUserID    int `gorm:"index"`
	ActualUtterance   string
	UtteredAt         time.Time
	UtteredAtTimezone string
}

func AddUtterance(actualUtterance string, telegramUserID int) (int64, error) {
	utterance := &Utterance{ActualUtterance: actualUtterance, TelegramUserID: telegramUserID, UtteredAt: time.Now().UTC(), UtteredAtTimezone: timezone}
	result := DB.Create(utterance)
	return utterance.ID, result.Error
}

func GetCount(atTime time.Time, telegramUserID int) int64 {
	loc, _ := time.LoadLocation(timezone)
	timeInZone := atTime.In(loc)
	startUTC, endUTC := getStartAndEndOfDayInLocation(timeInZone, loc)
	var count int64
	DB.Model(&Utterance{}).Where("uttered_at BETWEEN ? and ?", startUTC, endUTC).Count(&count)
	return count
}

func getStartAndEndOfDayInLocation(utcTime time.Time, loc *time.Location) (time.Time, time.Time) {
	timeInZone := utcTime.In(loc)
	y, m, d := timeInZone.Date()
	startOfDayInZone := time.Date(y, m, d, 0, 0, 0, 0, loc)
	endOfDayInZone := time.Date(y, m, d, 23, 59, 59, 59, loc)

	return startOfDayInZone.UTC(), endOfDayInZone.UTC()
}
