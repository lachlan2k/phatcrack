package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	log "github.com/sirupsen/logrus"
)

func AuditLog(c echo.Context, fields log.Fields, format string, args ...interface{}) {
	if fields == nil {
		fields = log.Fields{}
	}

	fields["log_type"] = "audit"
	fields["remote_ip"] = c.RealIP()

	user := auth.UserFromReq(c)
	if user == nil {
		fields["user_username"] = "Unknown"
		fields["user_id"] = "Unknown"
	} else {
		fields["user_username"] = user.Username
		fields["user_id"] = user.ID.String()
	}

	log.WithFields(fields).Warnf(format, args...)
}
