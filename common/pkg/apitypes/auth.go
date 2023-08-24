package apitypes

import "github.com/NHAS/webauthn/protocol"

type AuthLoginRequestDTO struct {
	Username string `json:"username" validate:"required,min=4,max=64,username"`
	Password string `json:"password" validate:"required,min=16,max=128"`
}

type AuthCurrentUserDTO struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

type AuthLoginResponseDTO struct {
	User                   AuthCurrentUserDTO `json:"user"`
	IsAwaitingMFA          bool               `json:"is_awaiting_mfa"`
	RequiresPasswordChange bool               `json:"requires_password_change"`
	RequiresMFAEnrollment  bool               `json:"requires_mfa_enrollment"`
}

type AuthWhoamiResponseDTO struct {
	User                   AuthCurrentUserDTO `json:"user"`
	IsAwaitingMFA          bool               `json:"is_awaiting_mfa"`
	RequiresPasswordChange bool               `json:"requires_password_change"`
	RequiresMFAEnrollment  bool               `json:"requires_mfa_enrollment"`
}

type AuthRefreshResponseDTO struct {
	User                   AuthCurrentUserDTO `json:"user"`
	IsAwaitingMFA          bool               `json:"is_awaiting_mfa"`
	RequiresPasswordChange bool               `json:"requires_password_change"`
	RequiresMFAEnrollment  bool               `json:"requires_mfa_enrollment"`
}

type AuthWebAuthnStartEnrollmentResponseDTO struct {
	protocol.CredentialCreation
}

type AuthWebAuthnStartChallengeResponseDTO struct {
	protocol.CredentialAssertion
}

type AuthChangePasswordRequestDTO struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=16,max=128"`
}
