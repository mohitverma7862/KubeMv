package auth

// Role is a KubeMv identity role. Phase 0 assigns a role and checks that a
// session exists. The administration phase enforces the role matrix. Kubernetes
// RBAC stays the authorization boundary once API calls exist.
type Role string

const (
	RoleSuperAdmin    Role = "SUPER_ADMIN"
	RolePlatformAdmin Role = "PLATFORM_ADMIN"
	RoleClusterAdmin  Role = "CLUSTER_ADMIN"
	RoleSRE           Role = "SRE"
	RoleDevOps        Role = "DEVOPS"
	RoleDeveloper     Role = "DEVELOPER"
	RoleSecurity      Role = "SECURITY"
	RoleViewer        Role = "VIEWER"
	RoleAuditor       Role = "AUDITOR"
)

// Catalog is the closed set of role names. Callers must not invent new ones.
func Catalog() []Role {
	return []Role{
		RoleSuperAdmin,
		RolePlatformAdmin,
		RoleClusterAdmin,
		RoleSRE,
		RoleDevOps,
		RoleDeveloper,
		RoleSecurity,
		RoleViewer,
		RoleAuditor,
	}
}

func KnownRole(role Role) bool {
	for _, candidate := range Catalog() {
		if candidate == role {
			return true
		}
	}
	return false
}
