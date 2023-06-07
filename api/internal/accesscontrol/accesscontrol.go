package accesscontrol

import (
	"fmt"

	"github.com/lachlan2k/phatcrack/api/internal/db"
)

func HasRightsToProject(user *db.User, project *db.Project) bool {
	if project.OwnerUserID == user.ID {
		return true
	}

	for _, share := range project.ProjectShare {
		if share.UserID == user.ID {
			return true
		}
	}

	return false
}

func HasRightsToProjectID(user *db.User, projId string) (bool, error) {
	proj, err := db.GetProjectForUser(projId, user.ID.String())
	if proj == nil || err == db.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get underlying project to check access control: %v", err)
	}
	return HasRightsToProject(user, proj), nil
}

func HasRightsToJobID(user *db.User, jobID string) (bool, error) {
	projId, err := db.GetJobProjID(jobID)
	if err == db.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get underlying job to check access control: %v", err)
	}
	return HasRightsToProjectID(user, projId)
}
