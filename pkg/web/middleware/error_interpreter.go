package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/ginEx"
)

func InterpretErrors(errorType gin.ErrorType, r errmap.ErrorRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
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
