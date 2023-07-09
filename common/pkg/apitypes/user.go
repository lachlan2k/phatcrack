package apitypes

type UserDTO struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

const UserRoleAdmin = "admin"
const UserRoleStandard = "standard"

var UserSignupRoles = []string{UserRoleStandard, UserRoleAdmin}
