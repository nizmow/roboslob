package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nizmow/roboslob/models"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Utterance{})

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		panic(err)
	}

	b.Handle(tb.OnText, func(m *tb.Message) {
		if matches_utterance(m.Text) {
			fmt.Println(m.Text)
		}
	})

	b.Start()
}

func matches_utterance(uttterance string) bool {
	utterances_to_match := [...]string{"foo", "bar"}

	for _, u := range utterances_to_match {
		if strings.Contains(uttterance, u) {
			return true
		}
	}
	return false
}
