package apitypes

type UserDTO struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

type UserMinimalDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type UsersGetAllResponseDTO struct {
	Users []UserMinimalDTO `json:"users"`
}

const UserRoleAdmin = "admin"
const UserRoleStandard = "standard"
const UserRoleServiceAccount = "service_account"
const UserRoleMFAExempt = "mfa_exempt"
const UserRoleRequiresPasswordChange = "requires_password_change"

var UserAssignableRoles = []string{UserRoleStandard, UserRoleAdmin, UserRoleServiceAccount, UserRoleMFAExempt, UserRoleRequiresPasswordChange}
