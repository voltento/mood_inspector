package main

import (
	"flag"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

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

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%v %v] %v", update.Message.From.FirstName, update.Message.From.LastName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, er := bot.Send(msg); er != nil {
			log.Println(er.Error())
		}
	}
}
