package dashboard_group

type DashboardConfig struct {
	ConfigId            string  `json:"configId,omitempty"`
	DashboardId         string  `json:"dashboardId,omitempty"`
	DescriptionOverride string  `json:"descriptionOverride,omitempty"`
	FiltersOverride     Filters `json:"filters,omitempty"`
	NameOverride        string  `json:"nameOverride,omitempty"`
}
