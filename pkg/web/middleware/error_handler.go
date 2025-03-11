// SPDX-FileCopyrightText: 2024 Greenbone Networks GmbH <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package middleware

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerFunc is a type for functions that handle errors.
// These functions return a boolean indicating whether they handled the error.
type ErrorHandlerFunc func(error, *gin.Context) bool

// ErrorHandlers is a slice of error handler functions.
// Each function is responsible for handling specific types of errors.
// All other error handling is currently done inside the handler for the specific endpoint.
var ErrorHandlers = []ErrorHandlerFunc{
	handleUnauthorizedError,
}

// ErrorHandler is a middleware that processes errors after the request context is finished.
// It iterates over registered error handlers to find one that can handle the error.
func ErrorHandler(errorType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Size() > 0 {
			return
		}

		for _, errorValue := range c.Errors.ByType(errorType) {
			for _, handler := range ErrorHandlers {
				if handler(errorValue.Err, c) {
					return
				}
			}
		}
	}
}

// handleUnauthorizedError handles errors related to unauthorized access by matching error messages with a regex pattern.
func handleUnauthorizedError(err error, c *gin.Context) bool {
	return handleRegexErrorHandler(err, c, `^could not bind header:.*$|^authorization failed:.*$`, http.StatusUnauthorized)
}

// handleRegexError matches an error message against a regex pattern and handles the error if it matches.
func handleRegexErrorHandler(err error, c *gin.Context, pattern string, statusCode int) bool {
	regex, _ := regexp.Compile(pattern)
	if regex.MatchString(err.Error()) {
		c.AbortWithStatus(statusCode)
		return true
	}
	return false
}
