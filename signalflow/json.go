package signalflow

import (
	"encoding/json"
	"fmt"
)

func parseJSONMessage(baseMessage Message, msg []byte) (Message, error) {
	var out Message
	switch baseMessage.Type() {
	case AuthenticatedType:
		out = &AuthenticatedMessage{}
	case ControlMessageType:
		var base BaseControlMessage
		if err := json.Unmarshal(msg, &base); err != nil {
			return nil, err
		}

		switch base.Event {
		case JobStartEvent:
			out = &JobStartControlMessage{}
		default:
			return &base, nil
		}
	case MetadataType:
		out = &MetadataMessage{}
	case MessageType:
		out = &MessageMessage{}
	default:
		return nil, fmt.Errorf("unknown message type: %s", baseMessage.Type())
	}
	err := json.Unmarshal(msg, out)
	return out, err
}
