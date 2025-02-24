package notification

type SplunkPlatformNotification struct {
	// Type sets which system to use to send the notification. For a Splunk Platform notification, this is always \"SplunkPlatform\".
	Type string `json:"type"`
	// Splunk Platform-supplied credential ID that Splunk Observability uses to authenticate the notification with the Splunk Platform system. Get this value from your Splunk Platform account settings.
	CredentialId string `json:"credentialId"`
}
