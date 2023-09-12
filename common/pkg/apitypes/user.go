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

var UserSignupRoles = []string{UserRoleStandard, UserRoleAdmin}
