package detector

// ReminderNotification struct for ReminderNotification
type ReminderNotification struct {
	Interval int64  `json:"interval"`
	Timeout  int64  `json:"timeout,omitempty"`
	Type     string `json:"type"`
}
