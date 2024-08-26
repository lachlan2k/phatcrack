package accesscontrol

import (
	"fmt"

	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/roles"
)

func HasOwnershipRightsToProject(user *db.User, project *db.Project) bool {
	return project.OwnerUserID == user.ID || user.HasRole(roles.UserRoleAdmin)
}

func HasRightsToProject(user *db.User, project *db.Project) bool {
	if project.OwnerUserID == user.ID {
		return true
	}

	if user.HasRole(roles.UserRoleAdmin) {
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
	if user.HasRole(roles.UserRoleAdmin) {
		return true, nil
	}

	proj, err := db.GetProjectForUser(projId, user)
	if proj == nil || err == db.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get underlying project to check access control: %w", err)
	}
	return HasRightsToProject(user, proj), nil
}

func HasRightsToJobID(user *db.User, jobID string) (bool, error) {
	projId, err := db.GetJobProjID(jobID)
	if err == db.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get underlying job to check access control: %w", err)
	}
	return HasRightsToProjectID(user, projId)
}
