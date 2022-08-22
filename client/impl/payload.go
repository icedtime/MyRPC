package impl

type SendMessage struct {
	Data string
}

type ReciveMessage struct {
	Data string
}

func NewSendMessage(msg string) *SendMessage {
	return &SendMessage{
		Data: msg,
	}
}

func NewReciveMessage(msg string) *ReciveMessage {
	return &ReciveMessage{
		Data: msg,
	}
}
