package pkg

import "sync"

type (
	ID int64
)

type Chats interface {
	AddChat(ID)
	Get() map[ID]struct{}
}

type DataBase struct {
	Chats Chats
}

func NewDatabase() DataBase {
	return DataBase{Chats: &chats{ids: make(map[ID]struct{})}}
}

type chats struct {
	ids map[ID]struct{}
	mtx sync.Mutex
}

func (ch *chats) AddChat(id ID) {
	ch.mtx.Lock()
	defer ch.mtx.Unlock()

	ch.ids[id] = struct{}{}
}

func (ch *chats) Get() map[ID]struct{} {
	ch.mtx.Lock()
	defer ch.mtx.Unlock()

	result := make(map[ID]struct{})
	for k, v := range ch.ids {
		result[k] = v
	}

	return result
}
