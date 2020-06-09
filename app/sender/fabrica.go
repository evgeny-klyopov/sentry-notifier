package sender

import (
	"errors"
)

func CreateObject(senderType string) (Sender, error) {
	var object Sender
	switch senderType {
	case TelegramType:
		object = &Telegram{}
	}

	if object == nil {
		return nil, errors.New("unsupported sender")
	}

	return object, nil
}
