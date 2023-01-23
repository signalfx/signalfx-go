package dashboard

// AclEntry - a single permission configuration entry
type AclEntry struct {
	// ID of the user, team, or organization for which you're granting permissions.
	PrincipalId string `json:"principalId"`
	// Clarify whether this permission configuration is for a user, a team, or an organization. Enum: [USER, TEAM, ORG]
	PrincipalType string `json:"principalType"`
	// Action the user, team, or organization can take with the dashboard group. Enum: [READ, WRITE]
	Actions []string `json:"actions"`
}
