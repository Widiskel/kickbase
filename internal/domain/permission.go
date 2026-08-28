package domain

// Permission constants for all domains (Create, Read, Update, Delete, Revert)
const (
	// Team Permissions
	PermTeamsCreate = "teams:create"
	PermTeamsRead   = "teams:read"
	PermTeamsUpdate = "teams:update"
	PermTeamsDelete = "teams:delete"
	PermTeamsRevert = "teams:revert"

	// Player Permissions
	PermPlayersCreate = "players:create"
	PermPlayersRead   = "players:read"
	PermPlayersUpdate = "players:update"
	PermPlayersDelete = "players:delete"
	PermPlayersRevert = "players:revert"

	// Match Permissions
	PermMatchesCreate = "matches:create"
	PermMatchesRead   = "matches:read"
	PermMatchesUpdate = "matches:update"
	PermMatchesDelete = "matches:delete"
	PermMatchesRevert = "matches:revert"

	// Result Permissions
	PermResultsCreate = "results:create"
	PermResultsRead   = "results:read"
	PermResultsUpdate = "results:update"
	PermResultsRevert = "results:revert"

	// Report Permissions
	PermReportsRead = "reports:read"

	// User/Admin Management Permissions
	PermUsersCreate = "users:create"
	PermUsersRead   = "users:read"
	PermUsersUpdate = "users:update"
	PermUsersDelete = "users:delete"
)

// RolePermissions maps a user role to its allowed permissions
var RolePermissions = map[string][]string{
	"admin": {
		PermTeamsCreate, PermTeamsRead, PermTeamsUpdate, PermTeamsDelete, PermTeamsRevert,
		PermPlayersCreate, PermPlayersRead, PermPlayersUpdate, PermPlayersDelete, PermPlayersRevert,
		PermMatchesCreate, PermMatchesRead, PermMatchesUpdate, PermMatchesDelete, PermMatchesRevert,
		PermResultsCreate, PermResultsRead, PermResultsUpdate, PermResultsRevert,
		PermReportsRead,
		PermUsersCreate, PermUsersRead, PermUsersUpdate, PermUsersDelete,
	},
	"staff": {
		PermTeamsRead, PermTeamsCreate, PermTeamsUpdate,
		PermPlayersRead, PermPlayersCreate, PermPlayersUpdate,
		PermMatchesRead, PermMatchesCreate, PermMatchesUpdate,
		PermResultsRead, PermResultsCreate,
		PermReportsRead,
	},
	"viewer": {
		PermTeamsRead,
		PermPlayersRead,
		PermMatchesRead,
		PermResultsRead,
		PermReportsRead,
	},
}

// GetPermissionsForRole returns all permissions granted to a given role
func GetPermissionsForRole(role string) []string {
	if perms, exists := RolePermissions[role]; exists {
		return perms
	}
	return []string{
		PermTeamsRead, PermPlayersRead, PermMatchesRead, PermResultsRead, PermReportsRead,
	}
}

// HasPermission checks whether a given role contains a specific permission
func HasPermission(role string, requiredPerm string) bool {
	perms := GetPermissionsForRole(role)
	for _, p := range perms {
		if p == requiredPerm || p == "*" {
			return true
		}
	}
	return false
}
