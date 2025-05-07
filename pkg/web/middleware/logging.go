// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package middleware

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/greenbone/opensight-golang-libraries/pkg/logs"
	"github.com/rs/zerolog"
)

const (
	// correlationIDHeader is the header key for the correlation ID
	correlationIDHeader = "X-Correlation-ID"
)

func Logging() gin.HandlerFunc {
	logHandler := logger.SetLogger(
		logger.WithLogger(func(c *gin.Context, _ zerolog.Logger) zerolog.Logger {
			correlationID := c.GetHeader(correlationIDHeader)
			if correlationID == "" {
				correlationID = uuid.New().String()
			}

			ctx := logs.WithCorrelationID(c.Request.Context(), correlationID)
			// update request context
			c.Request = c.Request.WithContext(ctx)

			// add correlation ID to response header
			c.Writer.Header().Set(correlationIDHeader, correlationID)

			// return context logger
			return *zerolog.Ctx(ctx)
		}),

		logger.WithUTC(true),
		// exclude alive endpoint to avoid log spam
		logger.WithSkipPath([]string{"/health/alive"}),
	)
	return logHandler
}
