package detector

// ReminderNotification struct for ReminderNotification
type ReminderNotification struct {
	// The interval at which you want to receive the notifications, in milliseconds.
	IntervalMs int64 `json:"interval"`
	// The duration during which repeat notifications are sent, in milliseconds.
	TimeoutMs int64 `json:"timeout,omitempty"`
	// Type of reminder notification. Currently, the only supported value is TIMEOUT.
	Type string `json:"type"`
}
