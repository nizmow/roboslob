package models

import "time"

type User struct {
	// db fields
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	TelegramID uint
	Username   string
	Timezone   string
}
