package ginEx

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestBindBody_BindingErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    string
		expectedMsg    string
		checkErrorType bool
	}{
		{
			name:           "Empty Body Returns BindingError",
			requestBody:    ``,
			expectedMsg:    "body can not be empty",
			checkErrorType: true,
		},
		{
			name:           "Unexpected EOF Returns BindingError",
			requestBody:    `{"name": "incomplete json`,
			expectedMsg:    "error parsing body",
			checkErrorType: true,
		},
		{
			name:        "Successful Bind without self-validation check",
			requestBody: `{"name": "expert"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.requestBody))

			var dto map[string]string
			result := BindBody(c, &dto)
			if result && !tt.checkErrorType {
				return
			}

			lastErr := c.Errors.Last().Err
			assert.Equal(t, tt.expectedMsg, lastErr.Error())

			var bindingError BindingError
			assert.True(t, errors.As(lastErr, &bindingError), "Error should be of type BindingError")
		})
	}
}

type sample struct {
}

func (s sample) Validate() models.ValidationErrors {
	return models.ValidationErrors{}
}

func TestBindBody_Validate(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewBufferString(`{
				"domain": "example.com",
				"port": 123,
				"isAuthenticationRequired": true,
				"isTlsEnforced": false,
				"username": "testUser",
				"password": "123"
			}`),
	)

	var s sample
	result := BindBody(c, &s)
	assert.False(t, result)

	lastErr := c.Errors.Last().Err
	assert.Contains(t, lastErr.Error(), "validation error")
}
