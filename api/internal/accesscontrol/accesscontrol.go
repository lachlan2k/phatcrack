package accesscontrol

import (
	"fmt"

	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
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

func HasRightsToProjectID(user *auth.UserClaims, projID string) (bool, error) {
	proj, err := db.GetProjectForUser(projID, user.ID)
	if proj == nil || err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get underlying project to check access control: %v", err)
	}
	return CanGetProject(user, proj), nil
}

func CanAccessJobID(user *auth.UserClaims, jobID string) (bool, error) {
	projID, err := db.GetJobProjID(jobID)
	if projID == "" || err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get underlying job to check access control: %v", err)
	}
	return HasRightsToProjectID(user, projID)
}
