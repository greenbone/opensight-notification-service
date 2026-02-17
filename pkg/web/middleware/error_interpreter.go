package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/ginEx"
)

func InterpretErrors(errorType gin.ErrorType, r errmap.ErrorRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.IsAborted() {
			return
		}

		for _, errorValue := range c.Errors.ByType(errorType) {
			actual := errorValue.Unwrap()

			bindingError := ginEx.BindingError{}
			if errors.As(actual, &bindingError) {
				c.AbortWithStatusJSON(http.StatusBadRequest, errorResponses.NewErrorValidationResponse("unable to parse the request", bindingError.Error(), nil))
				return
			}

			validationErrors := models.ValidationErrors{}
			if errors.As(actual, &validationErrors) {
				c.AbortWithStatusJSON(http.StatusBadRequest, errorResponses.NewErrorValidationResponse("", "", validationErrors))
				return
			}
			if errors.Is(actual, errs.ErrItemNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, errorResponses.NewErrorGenericResponse("item not found"))
				return
			}

			err, ok := r.Lookup(actual)
			if ok {
				c.AbortWithStatusJSON(err.Status, err.Response)
				return
			}

			// Unknown error
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponses.ErrorInternalResponse)
		}
	}
}
