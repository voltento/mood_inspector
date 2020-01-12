package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/voltento/mood_inspector/internal/pkg"
)

const (
	GREETING_MSG = "I'll be tracking you"
	REMIND_MSG   = ` 
1. What are you filling now?
2. Why are you filling this emotion?
3. What do you want now?
4. Why do you want it?
5. Which emotions do you need to get the goal?
6. How to transit yourself to these emotions?
`
)

type App struct {
	db       pkg.DataBase
	chatsMgr pkg.ChatsMgr
	bot      *tgbotapi.BotAPI
}

func main() {
	application := NewApp()
	application.Run()
}

func NewApp() *App {
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

	db := pkg.NewDatabase()
	chatsMgr := pkg.NewChatsMgr(db, bot)
	return &App{db: db, bot: bot, chatsMgr: chatsMgr}
}

func (app *App) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := app.bot.GetUpdatesChan(u)

	if err != nil {
		log.Fatal(err.Error())
	}

	timeoutChan := time.After(getReminderTimeout())

	for {
		select {
		case <-timeoutChan:
			timeoutChan = time.After(getReminderTimeout())
			if appropriateTime() {
				app.chatsMgr.SendMsgToAll(REMIND_MSG)
			}

		case update := <-updates:
			if update.Message != nil { // ignore any non-Message Updates
				app.db.Chats.AddChat(pkg.ID(update.Message.Chat.ID))

				log.Printf("[%v %v] %v", update.Message.From.FirstName, update.Message.From.LastName, update.Message.Text)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, GREETING_MSG)

				if _, er := app.bot.Send(msg); er != nil {
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

func getReminderTimeout() time.Duration {
	const timeout = time.Hour * 3
	floatTimeoutPart := time.Duration(uint64(time.Minute) * (rand.Uint64() % 60))
	return timeout + floatTimeoutPart
}
