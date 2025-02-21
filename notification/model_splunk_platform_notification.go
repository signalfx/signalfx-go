package notification

type SplunkPlatformNotification struct {
	// Type sets which system to use to send the notification. For a Splunk Platform notification, this is always \"SplunkPlatform\".
	Type string `json:"type"`
	// Tells SignalFX the HTTP Event Collector (HEC) URI for your Splunk platform instance
	Url string `json:"url"`
	//Enter the HTTP Event Collector token that allows access to your Splunk platform instance
	HecToken string `json:"hecToken"`
}
