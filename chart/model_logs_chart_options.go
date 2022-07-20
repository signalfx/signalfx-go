package chart

type Columns struct {
	Name string `json:"name,omitempty"`
}

type SortOptions struct {
	Field string `json:"field,omitempty"`
	Descending bool `json:"descending,omitempty"`
}