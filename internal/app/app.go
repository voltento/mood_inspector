package main

import (
	"flag"
	"log"
	"math/rand"
	"mood_inspector/internal/pkg"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type App struct {
	db pkg.DataBase
}

func main() {
	application := NewApp()
	application.Run()
}

func NewApp() *App {
	return &App{db: pkg.NewDatabase()}
}

func (app *App) Run() {
	secret := flag.String("secret", "", "pass secret to access api")
	flag.Parse()
	if secret == nil || *secret == "" {
		log.Fatal("expected secret")
	}

	bot, err := tgbotapi.NewBotAPI(*secret)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Fatal(err.Error())
	}

	const timeout = time.Hour * 3
	floatTimeoutPart := time.Duration(uint64(time.Minute) * (rand.Uint64() % 60))
	timeoutChan := time.After(timeout + floatTimeoutPart)

	for {
		select {
		case <-timeoutChan:
			timeoutChan = time.After(timeout)
			if appropriateTime() {
				for _, chatId := range app.db.Chats.Get() {
					msg := tgbotapi.NewMessage(int64(chatId), "What do you fill now?")

					if _, er := bot.Send(msg); er != nil {
						log.Println(er.Error())
					}
				}
			}

		case update := <-updates:
			if update.Message != nil { // ignore any non-Message Updates
				app.db.Chats.AddChat(pkg.ID(update.Message.Chat.ID))

				log.Printf("[%v %v] %v", update.Message.From.FirstName, update.Message.From.LastName, update.Message.Text)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I wil be tracking you")

				if _, er := bot.Send(msg); er != nil {
					log.Println(er.Error())
				}
			}
		}
	}
}

func appropriateTime() bool {
	now := time.Now().Hour()
	return now > 7 && now < 18
}
