package notification

import "time"

type Sender interface {
	Send(msg string)
}

type Notification interface {
	SendIfNeed(t time.Time, s Sender)
}

type TimeChecker interface {
	Check(time time.Time) bool
}

type MessageProvider interface {
	Message() string
}
