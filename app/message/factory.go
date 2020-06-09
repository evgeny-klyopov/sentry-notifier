package message

import (
	"errors"
)

func CreateObject(formatType string) (Messenger, error) {
	var object Messenger
	switch formatType {
	case MdFormatType:
		object = MD{}
	}

	if object == nil {
		return nil, errors.New("unsupported format message")
	}

	return object, nil
}
