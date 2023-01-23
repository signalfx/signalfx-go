/*
 * Integrations API for ServiceNow
 *
 * https://dev.splunk.com/observability/reference/api/integrations/latest
 * https://docs.splunk.com/Observability/admin/notif-services/servicenow.html
 */

package integration

// ServiceNowIntegration specifies the properties of a notification service integration between ServiceNow and Observability Cloud, in the form of a JSON object
type ServiceNowIntegration struct {
	// The creation date and time for the integration object, in Unix time UTC-relative. This value is "read-only".
	Created int64 `json:"created,omitempty"`
	// Observability-assigned user ID of the user that created the integration. If the system created the object, the value is "AAAAAAAAAA". This value is "read-only".
	Creator string `json:"creator,omitempty"`
	// Flag that indicates the state of the integration object. If `true`, the integration is enabled. If `false`, the integration is disabled, and you must enable it by setting "enabled" to `true` in a **PUT** request that updates the object. **NOTE:** Observability always sets the flag to `true` when you call  **POST** `/integration` to create an integration.
	Enabled bool `json:"enabled"`
	// Observability-assigned ID of an integration you create in the web UI or API. Use this property to retrieve an integration using the **GET**, **PUT**, or **DELETE** `/integration/{id}` endpoints or the **GET** `/integration/validate{id}/` endpoint, as described in this topic.
	Id string `json:"id,omitempty"`
	// The last time the integration was updated, in Unix time UTC-relative. This value is "read-only".
	LastUpdated int64 `json:"lastUpdated,omitempty"`
	// Observability-assigned ID of the last user who updated the integration. If the last update was by the system, the value is "AAAAAAAAAA". This value is "read-only".
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	// A human-readable label for the integration. This property helps you identify a specific integration when you're using multiple integrations for the same service.
	Name string `json:"name,omitempty"`
	// Integration type, always "ServiceNow"
	Type Type `json:"type"`
	// Name of the ServiceNow instance, for example `myInstances.service-now.com`. To learn more, see the documentation.
	InstanceName string `json:"instanceName,omitempty"`
	// Describe the type of issue, using standard **ITIL** terminology. The allowed values are:   * Incident   * Problem
	IssueType string `json:"issueType,omitempty"`
	// Username you created in ServiceNow for the Observability Cloud integration. **NOTE:** In ServiceNow, you have to assign the roles  `web_service_admin` and `itil` to this username.
	Username string `json:"username,omitempty"`
	// Password associated with the `username` you created for this integration.
	Password string `json:"password,omitempty"`
	// An optional template that Observability Cloud uses to create the ServiceNow POST JSON payloads when an alert sends a notification to ServiceNow. Use this optional field to send the values of Observability Cloud alert properties to specific fields in ServiceNow. See API reference for details.
	AlertTriggeredPayloadTemplate string `json:"alertTriggeredPayloadTemplate,omitempty"`
	// An optional template that Observability Cloud uses to create the ServiceNow PUT JSON payloads when an alert is cleared in ServiceNow. Use this optional field to send the values of Observability Cloud alert properties to specific fields in ServiceNow. See API reference for details.
	AlertResolvedPayloadTemplate string `json:"alertResolvedPayloadTemplate,omitempty"`
}
