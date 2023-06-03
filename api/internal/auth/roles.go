package auth

const RoleAdmin = "admin"
const RoleStandard = "standard"

// This is dynamic. We add this once the user completes MFA
const RoleMFACompleted = "mfa_completed"

const RoleMFAEnrolled = "mfa_enrolled"
const RoleMFAExempt = "mfa_exempt"
const RoleRequiresPasswordChange = "requires_password_change"

var RolesAllowedOnRegistration = []string{RoleAdmin, RoleStandard}

func AreRolesAllowedOnRegistration(roles []string) bool {
	for _, role := range roles {
		found := false
		for _, allowedRole := range roles {
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
