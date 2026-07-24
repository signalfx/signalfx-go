package notification

// Notification properties for an alert sent through a centralized email template.
type EmailTemplateNotification struct {
	// Tells SignalFx which system to use to send the notification. For an email template notification, this is always "EmailTemplate".
	Type string `json:"type"`
	// The SignalFx ID of the email template to render for this detector notification.
	TemplateId string `json:"templateId"`
}
