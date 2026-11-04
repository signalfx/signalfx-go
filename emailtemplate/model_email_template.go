package emailtemplate

// EmailTemplate represents a reusable detector alert email template.
type EmailTemplate struct {
	Id              string            `json:"id,omitempty"`
	Name            string            `json:"name,omitempty"`
	TriggerSubject  string            `json:"triggerSubject,omitempty"`
	TriggerBody     string            `json:"triggerBody,omitempty"`
	ResolvedSubject string            `json:"resolvedSubject,omitempty"`
	ResolvedBody    string            `json:"resolvedBody,omitempty"`
	To              []string          `json:"to,omitempty"`
	Cc              []string          `json:"cc,omitempty"`
	Bcc             []string          `json:"bcc,omitempty"`
	CustomHeaders   map[string]string `json:"customHeaders,omitempty"`
	CreatedOnMs     int64             `json:"createdOnMs,omitempty"`
	CreatedBy       string            `json:"createdBy,omitempty"`
	UpdatedOnMs     int64             `json:"updatedOnMs,omitempty"`
	UpdatedBy       string            `json:"updatedBy,omitempty"`
}

// SearchResult is the paginated email-template search response.
type SearchResult struct {
	Count   int              `json:"count,omitempty"`
	Results []*EmailTemplate `json:"results,omitempty"`
}
