package navigator

type TooltipKey struct {
	DisplayName       string `json:"displayName,omitempty"`
	Format            string `json:"format,omitempty"`
	IsSummaryProperty bool   `json:"isSummaryProperty"`
	Property          string `json:"property,omitempty"`
}
