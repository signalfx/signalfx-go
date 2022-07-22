package chart

type LogsListChartOptions struct {
	// The column configuration on logs chart
	Columns *Columns `json:"columns,omitempty"`
	// The sort option configuration for a logs chart
	SortOptions *SortOptions `json:"sortOptions,omitempty"`
	// The default connection configuration for a logs chart 
	DefaultConnection string `json:"defaultConnection,omitempty"`
}

type Columns struct {
	Name string `json:"name,omitempty"`
}

type SortOptions struct {
	Field string `json:"field,omitempty"`
	Descending bool `json:"descending,omitempty"`
}