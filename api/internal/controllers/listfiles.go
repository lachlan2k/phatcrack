package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/filerepo"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
	"github.com/lachlan2k/phatcrack/api/internal/roles"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookListsEndpoints(api *echo.Group) {
	api.GET("/all", handleGetAllListfiles)

	api.POST("/upload", handleListfileUpload)

	api.GET("/:id", handleGetListfile)
	api.DELETE("/:id", handleListfileDelete)
}

func handleListfileDelete(c echo.Context) error {
	id := c.Param("id")
	if !util.AreValidUUIDs(id) {
		return echo.ErrBadRequest
	}

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	AuditLog(c, log.Fields{
		"listfile_id": id,
	}, "User is deleting listfile")

	listfile, err := db.GetListfile(id)
	if err != nil {
		return util.ServerError("Failed to get listfile prior to deletion", err)
	}

	isAllowed := user.HasRole(roles.UserRoleAdmin) || listfile.CreatedByUserID == user.ID
	if !isAllowed && listfile.AttachedProjectID != nil {
		projId := listfile.AttachedProjectID.String()
		ok, err := accesscontrol.HasRightsToProjectID(user, projId)
		if err != nil {
			log.WithError(err).WithField("project_id", projId).WithField("user_id", user.ID.String()).Warn("Failed to check project access control for listfile")
		} else if ok {
			isAllowed = true
		}
	}

	if !isAllowed {
		return echo.ErrForbidden
	}

	err = db.MarkListfileForDeletion(id)
	if err != nil {
		return util.ServerError("Failed to mark listfile for deletion", err)
	}

	return c.JSON(http.StatusOK, "ok")
}

var validListfileTypes = []string{db.ListfileTypeRulefile, db.ListfileTypeWordlist}

func handleListfileUpload(c echo.Context) error {
	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	fileType := c.FormValue("file-type")
	if !slices.Contains(validListfileTypes, fileType) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid file type %s. Valid types are: %s", fileType, strings.Join(validListfileTypes, ", ")))
	}

	var projectID *uuid.UUID = nil
	projectIDStr := c.FormValue("project-id")

	if len(projectIDStr) > 0 {
		ok, err := accesscontrol.HasRightsToProjectID(user, projectIDStr)
		if err != nil {
			return util.GenericServerError(err)
		}
		if !ok {
			return echo.ErrForbidden
		}
		u, err := uuid.Parse(projectIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid project ID %s", projectIDStr))
		}
		projectID = &u
	}

	uploadedFile, err := c.FormFile("file")
	if err != nil {
		return util.ServerError("Failed to get file for upload. Perhaps disk space is low?", err)
	}

	filename := c.FormValue("file-name")
	if filename == "" {
		filename = uploadedFile.Filename
		if filename == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "No filename was given")
		}
	}

	maxFileSize := config.Get().General.MaximumUploadedFileSize
	if !user.HasRole(roles.UserRoleAdmin) && uploadedFile.Size > maxFileSize {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Uploaded file too large (maximum %d bytes)", maxFileSize))
	}

	maxAutodetectSize := config.Get().General.MaximumUploadedFileLineScanSize
	canAutodetectLineCount := uploadedFile.Size <= maxAutodetectSize

	lineCount, err := strconv.Atoi(c.FormValue("file-line-count"))
	if err != nil || (lineCount == 0 && !canAutodetectLineCount) || lineCount < 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("Invalid line count %d. Can only automatically detect lines for files up to %d bytes", lineCount, maxAutodetectSize),
		)
	}

	uploadedFileHandle, err := uploadedFile.Open()
	if err != nil {
		return util.ServerError("failed to open handle to file", err)
	}

	// Auto detect linecount
	if lineCount == 0 {
		lineCount = 1                // First line won't have a \n
		buf := make([]byte, 32*1024) // 32kb buffer
		lineSeparator := []byte{'\n'}

		for i := int64(0); i < maxAutodetectSize; i++ {
			n, err := uploadedFileHandle.Read(buf)
			lineCount += bytes.Count(buf[:n], lineSeparator)
			if err != nil {
				break
			}
		}

		// Rewind for when we need to copy it to disk
		uploadedFileHandle.Seek(0, io.SeekStart)
	}

	AuditLog(c, log.Fields{
		"listfile_size":      uploadedFile.Size,
		"listfile_linecount": lineCount,
		"listfile_filename":  filename,
		"listfile_type":      fileType,
	}, "User uploaded a new %s", fileType)

	// TODO rollback on later failures? We might run out of disk space on the io.Copy, etc
	listfile, err := db.CreateListfile(&db.Listfile{
		Name:              filename,
		FileType:          fileType,
		SizeInBytes:       uint64(uploadedFile.Size),
		Lines:             uint64(lineCount),
		CreatedByUserID:   user.ID,
		AttachedProjectID: projectID,
	})
	if err != nil {
		return util.ServerError("Failed to create new listfile", err)
	}

	outfile, err := filerepo.Create(listfile.ID)
	if err != nil {
		return util.ServerError("Failed to create new file on disk", err)
	}

	_, err = io.Copy(outfile, uploadedFileHandle)
	if err != nil {
		return util.ServerError("Failed to write to new file on disk", err)
	}

	err = db.MarkListfileAsAvailable(listfile.ID.String())
	if err != nil {
		return util.ServerError("Failed to prepare new listfile", err)
	}

	fleet.RequestFileDownload(listfile.ID)

	return c.JSON(http.StatusCreated, listfile.ToDTO())
}

func handleGetListfile(c echo.Context) error {
	id := c.Param("id")
	if !util.AreValidUUIDs(id) {
		return echo.ErrBadRequest
	}

	listfile, err := db.GetListfile(id)
	if err == db.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return util.ServerError("Failed to fetch wordlist", err)
	}

	if listfile.AttachedProjectID != nil {
		user := auth.UserFromReq(c)
		if user == nil {
			return echo.ErrForbidden
		}

		ok, err := accesscontrol.HasRightsToProjectID(user, listfile.AttachedProjectID.String())
		if err != nil {
			log.WithError(err).
				WithField("project_id", listfile.AttachedProjectID.String()).
				WithField("user_id", user.ID.String()).
				Warn("Failed to lookup user access to listifle")

			return echo.ErrNotFound
		}

		if !ok {
			return echo.ErrNotFound
		}
	}

	return c.JSON(http.StatusOK, listfile.ToDTO())
}

func handleGetAllListfiles(c echo.Context) error {
	listfiles, err := db.GetAllPublicListfiles()
	if err != nil {
		return util.ServerError("Failed to fetch listfiles", err)
	}

	var res apitypes.GetAllListfilesDTO
	res.Listfiles = make([]apitypes.ListfileDTO, len(listfiles))
	for i, lf := range listfiles {
		res.Listfiles[i] = lf.ToDTO()
	}

	return c.JSON(http.StatusOK, res)
}
