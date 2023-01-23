package dashboard

// ObjectPermissions - available if your organization has the \"access_control\" feature enabled. Read and write permission configuration to specify which user, team, and organization can view and/or edit your dashboard group.
type ObjectPermissions struct {
	// The ID of the dashboard group you want your dashboard to inherit permissions from.
	Parent string `json:"parent,omitempty"`
	// List of custom permission configurations if you don't inherit permissions
	Acl []*AclEntry `json:"acl,omitempty"`
}
