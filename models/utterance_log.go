package models

import "time"

type Utterance struct {
	// db fields
	ID        uint
	UpdatedAt time.Time
	CreatedAt time.Time

	UserID            uint `gorm:"index"`
	UtteredAt         time.Time
	UtteredAtTimezone string
}
