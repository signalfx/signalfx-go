package integration

type SplunkPlatformIntegration struct {
	// The creation date and time for the integration object, in Unix time UTC-relative. The system sets this value, and you can't modify it.
	Created int64 `json:"created,omitempty"`
	// Splunk Observability assigned user ID of the user that created the integration object. If the system created the object, the value is \"AAAAAAAAAA\". The system sets this value, and you can't modify it.
	Creator string `json:"creator,omitempty"`
	//Name of the user that created the integration
	CreatedByName string `json:"createdByName,omitempty"`
	//Type of service that this integration represents, in the form of an enumerated string, always "SplunkPlatform"
	Type Type `json:"type"`
	// Flag that indicates the state of the integration object. If  `true`, the integration is enabled. If `false`, the integration is disabled, and you must enable it by setting \"enabled\" to `true` in a **PUT** request that updates the object. <br> **NOTE:** Splunk Observability always sets the flag to `true` when you call  **POST** `/integration` to create an integration.
	Enabled bool `json:"enabled"`
	//HTTP Event Collector token that allows access to your Splunk platform instance
	HecToken string `json:"hecToken,omitempty"`
	//Splunk Observability assignedID of an integration you create in the web UI or API. Use this property to retrieve an integration using the **GET**, **PUT**, or **DELETE** `/integration/{id}` endpoints or the **GET** `/integration/validate{id}/` endpoint, as described in this topic.
	Id string `json:"id,omitempty"`
	// The last time the integration was updated, in Unix time UTC-relative. This value is \"read-only\".
	LastUpdated int64 `json:"lastUpdated,omitempty"`
	//Splunk Observability assigned ID of the last user who updated the integration. If the last update was by the system, the value is \"AAAAAAAAAA\". This value is \"read-only\".
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	//Name of the user that last updated the integration
	LastUpdatedByName string `json:"lastUpdatedByName,omitempty"`
	// A human-readable label for the integration. This property helps you identify a specific integration when you're using multiple integrations for the same service.
	Name string `json:"name,omitempty"`
	//Customize the Splunk platform alert payload using Handlebars syntax
	PayloadTemplate string `json:"payloadTemplate,omitempty"`
	//Specify the HTTP Event Collector (HEC) URI for your Splunk platform instance
	Url string `json:"url,omitempty"`
}
