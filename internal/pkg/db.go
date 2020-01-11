package pkg

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
}

func (ch *chats) AddChat(id ID) {
	ch.ids[id] = struct{}{}
}

func (ch *chats) Get() []ID {
	var ids []ID
	for id := range ch.ids {
		ids = append(ids, id)
	}
	return ids
}