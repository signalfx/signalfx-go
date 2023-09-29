package event

type Event struct {
	ID                 string                 `json:"id"`
	Metadata           map[string]interface{} `json:"metadata"`
	Properties         map[string]interface{} `json:"properties"`
	SFEventCategory    EventCategory          `json:"sf_eventCategory"`
	SFEventCreatedOnMs int64                  `json:"sf_eventCreatedOnMs"`
	SFEventType        string                 `json:"sf_eventType"`
	Timestamp          int64                  `json:"timestamp"`
	TsID               string                 `json:"tsId"`
}

type EventCategory string

const (
	UserDefined      EventCategory = "USER_DEFINED"
	Alert            EventCategory = "ALERT"
	Audit            EventCategory = "AUDIT"
	Job              EventCategory = "JOB"
	Collectd         EventCategory = "COLLECTD"
	ServiceDiscovery EventCategory = "SERVICE_DISCOVERY"
	Exception        EventCategory = "EXCEPTION"
)

type Events []Event
