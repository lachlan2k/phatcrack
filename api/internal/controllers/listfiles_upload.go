package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/filerepo"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
	"github.com/lachlan2k/phatcrack/api/internal/roles"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	log "github.com/sirupsen/logrus"
)

var validListfileTypes = []string{db.ListfileTypeRulefile, db.ListfileTypeWordlist}

type listfileUploadForm struct {
	fileType         string
	fileName         string
	filePartFileName string
	fileSize         int64
	fileSeen         bool

	lineCount *int

	projectID *uuid.UUID
}

func handleListfileUpload(c echo.Context) error {
	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	tmpFile, tmpFilePath, err := filerepo.MakeTmp()
	if err != nil {
		return util.ServerError("Failed to create temporary file", err)
	}

	success := false
	defer func() {
		if !success {
			os.Remove(tmpFilePath)
		}
	}()

	maxFileSize := config.Get().General.MaximumUploadedFileSize
	if user.HasRole(roles.UserRoleAdmin) {
		maxFileSize = 0
	}

	mpReader, err := c.Request().MultipartReader()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid multipart request")
	}

	form, err := parseListfileUploadForm(mpReader, tmpFile, maxFileSize)
	if err != nil {
		return err
	}

	if !form.fileSeen {
		return echo.NewHTTPError(http.StatusBadRequest, "No file was uploaded")
	}

	if !slices.Contains(validListfileTypes, form.fileType) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid file type %q. Valid types are: %s", form.fileType, strings.Join(validListfileTypes, ", ")))
	}

	if form.projectID != nil {
		ok, err := accesscontrol.HasRightsToProjectID(user, form.projectID.String())
		if err != nil {
			return util.GenericServerError(err)
		}
		if !ok {
			return echo.ErrForbidden
		}
	}

	if form.fileName == "" {
		form.fileName = form.filePartFileName
		if form.fileName == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "No filename was given")
		}
	}
	if len(form.fileName) > 255 {
		return echo.NewHTTPError(http.StatusBadRequest, "Filename too long")
	}

	maxAutodetectSize := config.Get().General.MaximumUploadedFileLineScanSize
	canAutodetectLineCount := form.fileSize <= maxAutodetectSize

	requestedAutoDetectLineCount := form.lineCount == nil || *form.lineCount == 0

	if requestedAutoDetectLineCount && !canAutodetectLineCount {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot autodetect line count, file too large")
	}

	if requestedAutoDetectLineCount {
		n, err := detectFileLineCount(tmpFile)
		if err != nil {
			return util.ServerError("Failed to autodetect line count", err)
		}
		form.lineCount = &n
	}

	AuditLog(c, log.Fields{
		"listfile_size":      form.fileSize,
		"listfile_linecount": *form.lineCount,
		"listfile_filename":  form.fileName,
		"listfile_type":      form.fileType,
	}, "User uploaded a new %s", form.fileType)

	listfile, err := db.CreateListfile(&db.Listfile{
		Name:              form.fileName,
		FileType:          form.fileType,
		SizeInBytes:       uint64(form.fileSize),
		Lines:             uint64(*form.lineCount),
		CreatedByUserID:   user.ID,
		AttachedProjectID: form.projectID,
	})
	if err != nil {
		return util.ServerError("Failed to create new listfile", err)
	}

	success = true
	err = filerepo.CreateFromTmp(listfile.ID, tmpFilePath)
	if err != nil {
		return util.ServerError("Failed to move file to disk", err)
	}

	err = db.MarkListfileAsAvailable(listfile.ID.String())
	if err != nil {
		return util.ServerError("Failed to prepare new listfile", err)
	}

	fleet.RequestFileDownload(listfile.ID)

	return c.JSON(http.StatusCreated, listfile.ToDTO())
}

func parseListfileUploadForm(mpReader *multipart.Reader, tmpFile *os.File, maxFileSize int64) (*listfileUploadForm, error) {
	f := &listfileUploadForm{}

	for {
		part, err := mpReader.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to read multipart request")
		}

		switch part.FormName() {
		case "file-type":
			if f.fileType != "" {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "File type already set")
			}

			ft, err := io.ReadAll(part)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to read file type")
			}
			f.fileType = string(ft)

		case "file-name":
			if f.fileName != "" {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "File name already set")
			}

			fn, err := io.ReadAll(part)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to read file name")
			}
			f.fileName = string(fn)

		case "file-line-count":
			if f.lineCount != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "Line count already set")
			}

			lc, err := io.ReadAll(part)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to read line count")
			}
			lcI, err := strconv.Atoi(string(lc))
			f.lineCount = &lcI
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse line count")
			}

		case "project-id":
			if f.projectID != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "Project ID already set")
			}

			pid, err := io.ReadAll(part)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to read project ID")
			}
			pidParsed, err := uuid.Parse(string(pid))
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse project ID")
			}
			f.projectID = &pidParsed

		case "file":
			if f.fileSeen {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "File already set")
			}
			f.fileSeen = true
			f.filePartFileName = part.FileName()

			if maxFileSize > 0 {
				n, err := io.CopyN(tmpFile, part, maxFileSize+1)
				if err != nil && !errors.Is(err, io.EOF) {
					return nil, util.ServerError("Failed to upload file. Perhaps disk space is low?", err)
				}

				if n == maxFileSize+1 {
					return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("File too large (maximum %d bytes)", maxFileSize))
				}

				f.fileSize = n
			} else {
				n, err := io.Copy(tmpFile, part)
				if err != nil && !errors.Is(err, io.EOF) {
					return nil, util.ServerError("Failed to upload file. Perhaps disk space is low?", err)
				}

				f.fileSize = n
			}

		default:
			return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unexpected form field %s", part.FormName()))
		}

		part.Close()
	}

	return f, nil
}

func detectFileLineCount(f *os.File) (int, error) {
	lineCount := 1

	buf := make([]byte, 32*1024) // 32kb buffer
	f.Seek(0, io.SeekStart)
	lineSeparator := []byte{'\n'}

	for {
		n, err := f.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return 0, err
		}

		lineCount += bytes.Count(buf[:n], lineSeparator)
		if err != nil {
			break
		}
	}

	return lineCount, nil
}
