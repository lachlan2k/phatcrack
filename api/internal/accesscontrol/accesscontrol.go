package accesscontrol

import (
	"fmt"

	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
)

func HasRightsToProject(user *auth.UserClaims, project *db.Project) bool {
	if project.OwnerUserID.String() == user.ID {
		return true
	}

	for _, share := range project.ProjectShare {
		if share.UserID.String() == user.ID {
			return true
		}
	}

	return false
}

func HasRightsToProjectID(user *auth.UserClaims, projId string) (bool, error) {
	proj, err := db.GetProjectForUser(projId, user.ID)
	if proj == nil || err == db.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get underlying project to check access control: %v", err)
	}
	return HasRightsToProject(user, proj), nil
}

func HasRightsToJobID(user *auth.UserClaims, jobID string) (bool, error) {
	projId, err := db.GetJobProjID(jobID)
	if err == db.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get underlying job to check access control: %v", err)
	}
	return HasRightsToProjectID(user, projId)
}
