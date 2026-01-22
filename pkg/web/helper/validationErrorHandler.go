package helper

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
)

type ValidateErrors map[string]string

func (v ValidateErrors) Error() string {
	return "validation error"
}

func ValidationErrorHandler(errorType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, errorValue := range c.Errors.ByType(errorType) {
			validateErrors := ValidateErrors{}
			if errors.As(errorValue, &validateErrors) {
				c.AbortWithStatusJSON(http.StatusBadRequest, errorResponses.NewErrorValidationResponse("", "", validateErrors))
				return
			}
		}
	}
}
