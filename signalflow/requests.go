package signalflow

type AuthType string

func (at AuthType) MarshalJSON() ([]byte, error) {
	return []byte(`"authenticate"`), nil
}

type AuthRequest struct {
	Type AuthType `json:"type"`
	// The Auth token for the org
	Token string `json:"token"`
}

type ExecuteType string

func (et ExecuteType) MarshalJSON() ([]byte, error) {
	return []byte(`"execute"`), nil
}

type ExecuteRequest struct {
	Type       ExecuteType `json:"type"`
	Program    string      `json:"program"`
	Channel    string      `json:"channel"`
	Start      int64       `json:"start"`
	Stop       int64       `json:"stop"`
	Resolution int32       `json:"resolution"`
	MaxDelay   int32       `json:"maxDelay"`
	Immediate  bool        `json:"immediate"`
	Timezone   string      `json:"timezone"`
}
