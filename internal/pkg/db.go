package pkg

import "sync"

type (
	ID int64
)

type Chats interface {
	AddChat(ID)
	Get() []ID
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

func (ch *chats) Get() []ID {
	ch.mtx.Lock()
	defer ch.mtx.Unlock()

	ids := make([]ID, 0, 0)
	for id := range ch.ids {
		ids = append(ids, id)
	}
	return ids
}
