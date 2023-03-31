package accesscontrol

import (
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
)

func CanGetProject(user *auth.UserClaims, project *db.Project) bool {
	if project.OwnerUserID.Hex() == user.ID {
		return true
	}

	for _, id := range project.SharedWithUserIDs {
		if id.Hex() == user.ID {
			return true
		}
	}

	return false
}
