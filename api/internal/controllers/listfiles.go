package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/filerepo"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookListsEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong lists")
	})

	api.POST("/upload", handleListfileUpload)

	api.GET("/wordlist/all", handleGetAllWordlists)
	api.GET("/rulefile/all", handleGetAllRuleFiles)

	api.GET("/wordlist/:id", handleGetWordlist)
	api.GET("/rulefile/:id", handlGetRuleFile)
}

func handleListfileUpload(c echo.Context) error {
	fileType := c.FormValue("file-type")
	if fileType != db.ListfileTypeRulefile && fileType != db.ListfileTypeWordlist {
		return echo.ErrBadRequest
	}

	uploadedFile, err := c.FormFile("file")
	if err != nil {
		return echo.ErrBadRequest
	}

	maxFileSize := config.Get().MaximumUploadedFileSize
	if uploadedFile.Size > maxFileSize {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Uploaded file too large (maximum %d bytes)", maxFileSize))
	}

	maxAutodetectSize := config.Get().MaximumUploadedFileLineScanSize
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

	// TODO rollback on later failures
	listfile, err := db.CreateListfile(&db.Listfile{
		Name:        uploadedFile.Filename,
		FileType:    fileType,
		SizeInBytes: uint64(uploadedFile.Size),
		Lines:       uint64(lineCount),
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

func handleGetWordlist(c echo.Context) error {
	id := c.Param("id")
	if !util.AreValidUUIDs(id) {
		return echo.ErrBadRequest
	}

	list, err := db.GetListfile(id)
	if err == db.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return util.ServerError("Failed to fetch wordlist", err)
	}

	return c.JSON(http.StatusOK, list.ToDTO())
}

func handlGetRuleFile(c echo.Context) error {
	id := c.Param("id")
	if !util.AreValidUUIDs(id) {
		return echo.ErrBadRequest
	}

	list, err := db.GetListfile(id)
	if err == db.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return util.ServerError("Failed to fetch rulefile", err)
	}

	return c.JSON(http.StatusOK, list.ToDTO())
}

func handleGetAllWordlists(c echo.Context) error {
	lists, err := db.GetAllWordlists()
	if err != nil {
		return util.ServerError("Failed to fetch wordlists", err)
	}

	var res apitypes.GetAllWordlistsDTO
	res.Wordlists = make([]apitypes.ListfileDTO, len(lists))
	for i, list := range lists {
		res.Wordlists[i] = list.ToDTO()
	}

	return c.JSON(http.StatusOK, res)
}

func handleGetAllRuleFiles(c echo.Context) error {
	lists, err := db.GetAllRulefiles()
	if err != nil {
		return util.ServerError("Failed to fetch rulefiles", err)
	}

	var res apitypes.GetAllRuleFilesDTO
	res.RuleFiles = make([]apitypes.ListfileDTO, len(lists))
	for i, list := range lists {
		res.RuleFiles[i] = list.ToDTO()
	}

	return c.JSON(http.StatusOK, res)
}
