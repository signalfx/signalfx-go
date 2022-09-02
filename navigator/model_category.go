package navigator

type Category struct {
	CategoryCode          string `json:"categoryCode,omitempty"`
	CategoryDisplayName   string `json:"categoryDisplayName,omitempty"`
	CategoryGroupName     string `json:"categoryGroupDisplayName,omitempty"`
	CategoryInstanceLabel string `json:"categoryInstanceLabel,omitempty"`
	ConnectedCategory     bool   `json:"connectedCategory"`
}
