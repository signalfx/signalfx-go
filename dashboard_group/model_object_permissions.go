package dashboard_group

// ObjectPermissions - available if your organization has the \"access_control\" feature enabled. Read and write permission configuration to specify which user, team, and organization can view and/or edit your dashboard group.
type ObjectPermissions struct {
	// The parent property is null by default for dashboard groups and can't be changed
	Parent string `json:"parent,omitempty"`
	// List of permission configurations
	Acl []*AclEntry `json:"acl,omitempty"`
}
