package middleware

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
)

func InterpretErrors(errorType gin.ErrorType, r errmap.ErrorRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, errorValue := range c.Errors.ByType(errorType) {

			actual := errorValue.Unwrap()

			if isBindingError(actual) {
				c.AbortWithStatusJSON(http.StatusBadRequest, errorResponses.NewErrorValidationResponse("unable to parse the request", errorValue.Error(), nil))
				return
			}

			validationErrors := models.ValidationErrors{}
			if errors.As(actual, &validationErrors) {
				c.AbortWithStatusJSON(http.StatusBadRequest, errorResponses.NewErrorValidationResponse("", "", validationErrors))
				return
			}

			err, ok := r.Lookup(actual)
			if ok {
				c.AbortWithStatusJSON(err.Status, err.Response)
				return
			}

			// Unknown error
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponses.ErrorInternalResponse)
			return
		}
	}
}

func isBindingError(err error) bool {
	var syntaxErr *json.SyntaxError
	var unmarshalErr *json.UnmarshalTypeError

	return errors.As(err, &syntaxErr) ||
		errors.As(err, &unmarshalErr) ||
		errors.Is(err, io.EOF) ||
		errors.Is(err, io.ErrUnexpectedEOF)
}
