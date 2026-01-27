package ginEx

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/models"
)

type BindingError struct {
	message string
}

func newBindingError(message string) BindingError {
	return BindingError{message: message}
}

func (e BindingError) Error() string {
	return e.message
}

// Validate can be implemented by a dto that is used in BindAndValidateBody for custom validation
type Validate interface {
	Validate() models.ValidationErrors
}

func BindAndValidateBody(c *gin.Context, bodyDto any) bool {
	err := c.ShouldBindJSON(bodyDto)
	if err != nil {
		if errors.Is(err, io.EOF) {
			_ = c.Error(newBindingError("body can not be empty"))
			return false
		} else if errors.Is(err, io.ErrUnexpectedEOF) {
			_ = c.Error(newBindingError("error parsing body"))
			return false
		}

		if isUnmarshallError(err) {
			_ = c.Error(newBindingError("error unmarshalling body"))
			return false
		}

		_ = c.Error(err)
		return false
	}

	if value, ok := bodyDto.(Validate); ok {
		err := value.Validate()
		if err != nil || len(err) == 0 {
			_ = c.Error(err)
			return false
		}
	}

	return !c.IsAborted()
}

func AddError(c *gin.Context, err error) (errorIsNotNil bool) {
	if err == nil {
		return false
	}

	_ = c.Error(err)
	return true
}

func isUnmarshallError(err error) bool {
	var syntaxErr *json.SyntaxError
	var unmarshalErr *json.UnmarshalTypeError
	var invalidUnmarshalErr *json.InvalidUnmarshalError

	return errors.As(err, &syntaxErr) || errors.As(err, &unmarshalErr) || errors.As(err, &invalidUnmarshalErr)
}
