package main

import (
	"encoding/json"
	"fmt"
	"gitlab.mobbtech.com/iqbus/iqbus_go_client/errors"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/voltento/mood_inspector/internal/pkg"
	"github.com/voltento/mood_inspector/internal/pkg/notification"
)

const (
	GREETING_MSG = "!I'll be tracking you"
	REMIND_MSG   = `1. What are you filling now?
2. Why are you filling this emotion?
3. What do you want now?
4. Why do you want it?
5. Which emotions do you need to get the goal?
6. How to transit yourself to these emotions?`
)

type App struct {
	db            pkg.DataBase
	chatsMgr      pkg.ChatsMgr
	bot           *tgbotapi.BotAPI
	notifications []notification.Notification
	sender        notification.Sender
}

func main() {
	cfg := pkg.BuildConfig()
	application := NewApp(cfg)
	application.Run()
}

func NewApp(config *pkg.Config) *App {
	if data, err := json.Marshal(config); err != nil {
		log.Fatalf("can not marshal config: %v", err.Error())
	} else {
		log.Printf("loaded config: %v", string(data))
	}

	if len(config.TelegramBot.Id) == 0 {
		log.Fatalf("expected not empty value for config section telegram_bot.id")
	}

	if len(config.TelegramBot.Secret) == 0 {
		log.Fatalf("expected not empty value for config section telegram_bot.secret")
	}

	bot, err := tgbotapi.NewBotAPI(fmt.Sprintf("%v:%v", config.TelegramBot.Id, config.TelegramBot.Secret))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	bot.Self = tgbotapi.User{ID: 867251407, IsBot: true}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	db := pkg.NewDatabase()
	chatsMgr := pkg.NewChatsMgr(db, bot)
	notifications, err := NewNotifications(config)
	if err != nil {
		log.Panic(err.Error())
	}

	return &App{
		db:            db,
		bot:           bot,
		chatsMgr:      chatsMgr,
		notifications: notifications,
		sender:        &sender{chatsMgr},
	}
}

func NewNotifications(config *pkg.Config) ([]notification.Notification, error) {
	result := make([]notification.Notification, 0, len(config.Notifications))
	for _, notificationCfg := range config.Notifications {
		n, err := notification.NewNotification(&notificationCfg)
		if err != nil {
			return nil, errors.Wrap(err, "can not creat notification from %v")
		}
		result = append(result, n)
	}
	return result, nil
}

func (a *App) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := a.bot.GetUpdatesChan(u)

	if err != nil {
		log.Fatal(err.Error())
	}

	const tick = time.Second

	for {
		select {
		case <-time.After(tick):
			a.ProcessSendTick()

		case update := <-updates:
			if update.Message != nil { // ignore any non-Message Updates
				a.db.Chats.AddChat(pkg.ID(update.Message.Chat.ID))

				log.Printf("[%v %v] %v", update.Message.From.FirstName, update.Message.From.LastName, update.Message.Text)

				text := fmt.Sprintf("%v\nYou will recieve this messages several times a day: \n%v", GREETING_MSG, REMIND_MSG)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

				if _, er := a.bot.Send(msg); er != nil {
					log.Println(er.Error())
				}
			}
		}
	}
}

func (a *App) ProcessSendTick() {
	now := time.Now()
	for _, n := range a.notifications {
		n.SendIfNeed(now, a.sender)
	}
}

type sender struct {
	mgr pkg.ChatsMgr
}

func (s *sender) Send(msg string) {
	s.mgr.SendMsgToAll(msg)
}
