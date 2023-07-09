package auth

const RoleAdmin = "admin"
const RoleStandard = "standard"

const RoleMFAEnrolled = "mfa_enrolled"
const RoleMFAExempt = "mfa_exempt"
const RoleRequiresPasswordChange = "requires_password_change"

var RolesAllowedOnRegistration = []string{RoleAdmin, RoleStandard}

func AreRolesAllowedOnRegistration(roles []string) bool {
	for _, role := range roles {
		found := false
		for _, allowedRole := range RolesAllowedOnRegistration {
			if role == allowedRole {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}
