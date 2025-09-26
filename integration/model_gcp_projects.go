package integration

type GCPProjects struct {
	SelectedProjectIds []string `json:"selectedProjectIds,omitempty"`
	SyncMode           SyncMode `json:"syncMode,omitempty"`
}

type SyncMode string

const (
	SELECTED      SyncMode = "SELECTED"
	ALL_REACHABLE SyncMode = "ALL_REACHABLE"
)
