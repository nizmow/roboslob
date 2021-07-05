package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nizmow/roboslob/models"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		panic("failed to connect database")
	}
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Utterance{})
	models.SetDB(db)

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		panic(err)
	}

	b.Handle("/count", func(m *tb.Message) {
		count := models.GetCount(time.Now().UTC(), m.Sender.ID)
		b.Send(m.Chat, fmt.Sprintf("%s has count %d", m.Sender.Username, count))
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		if matches_utterance(m.Text) {
			models.AddUtterance(m.Text, m.Sender.ID)
		}
	})

	b.Start()
}

func matches_utterance(uttterance string) bool {
	utterances_to_match := [...]string{"🥝🎂", "🥝🍰"}

	for _, u := range utterances_to_match {
		if strings.Contains(uttterance, u) {
			return true
		}
	}
	return false
}
