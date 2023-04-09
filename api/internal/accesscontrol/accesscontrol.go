package accesscontrol

import (
	"fmt"

	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/dbnew"
)

func HasRightsToProject(user *auth.UserClaims, project *dbnew.Project) bool {
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
	proj, err := dbnew.GetProjectForUser(projId, user.ID)
	if proj == nil || err == dbnew.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get underlying project to check access control: %v", err)
	}
	return HasRightsToProject(user, proj), nil
}

func HasRightsToJobID(user *auth.UserClaims, jobID string) (bool, error) {
	projId, err := dbnew.GetJobProjID(jobID)
	if err == dbnew.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get underlying job to check access control: %v", err)
	}
	return HasRightsToProjectID(user, projId)
}
