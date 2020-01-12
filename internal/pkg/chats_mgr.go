package pkg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type ChatsMgr interface {
	SendMsgToAll(msg string)
}

type chatsMgr struct {
	db  DataBase
	bot *tgbotapi.BotAPI
}

func NewChatsMgr(db DataBase, bot *tgbotapi.BotAPI) ChatsMgr {
	return &chatsMgr{db: db, bot: bot}
}

func (chMgr *chatsMgr) SendMsgToAll(msg string) {
	for _, chatId := range chMgr.db.Chats.Get() {
		msg := tgbotapi.NewMessage(int64(chatId), msg)

		if _, er := chMgr.bot.Send(msg); er != nil {
			log.Println(er.Error())
		}
	}
}
