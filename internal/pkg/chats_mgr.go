package pkg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type ChatsMgr interface {
	SendMsgToAll(msg string)
	SendMsg(msg string, chat ID)
}

type chatsMgr struct {
	db  DataBase
	bot *tgbotapi.BotAPI
}

func NewChatsMgr(db DataBase, bot *tgbotapi.BotAPI) ChatsMgr {
	return &chatsMgr{db: db, bot: bot}
}

func (chMgr *chatsMgr) SendMsg(msg string, chatId ID) {
	if _, isOk := chMgr.db.Chats.Get()[chatId]; !isOk {
		return
	}

	chMgr.sendMsgImpl(chatId, msg)
}

func (chMgr *chatsMgr) SendMsgToAll(msg string) {
	for chatId := range chMgr.db.Chats.Get() {
		chMgr.sendMsgImpl(chatId, msg)
	}
}

func (chMgr *chatsMgr) sendMsgImpl(chatId ID, msg string) {
	tgMessage := tgbotapi.NewMessage(int64(chatId), msg)

	if _, er := chMgr.bot.Send(tgMessage); er != nil {
		log.Println(er.Error())
	}
}
