package roles

const UserRoleAdmin = "admin"
const UserRoleStandard = "standard"
const UserRoleServiceAccount = "service_account"
const UserRoleMFAExempt = "mfa_exempt"
const UserRoleRequiresPasswordChange = "requires_password_change"

const UserRoleMFAEnrolled = "mfa_enrolled"

var UserAssignableRoles = []string{UserRoleStandard, UserRoleAdmin, UserRoleServiceAccount, UserRoleMFAExempt, UserRoleRequiresPasswordChange}

func AreRolesAssignable(roles []string) bool {
	for _, role := range roles {
		found := false
		for _, allowedRole := range UserAssignableRoles {
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
