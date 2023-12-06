package util

func ValidatePasswordStrength(password string) (ok bool, feedback string) {
	if len(password) < 16 {
		return false, "Password should be a 16 characters minimum"
	}
	return true, ""
}
