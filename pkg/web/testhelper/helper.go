package testhelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/greenbone/keycloak-client-golang/auth"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	"github.com/stretchr/testify/assert"
)

func VerifyResponseWithMetadata[T any](
	t *testing.T,
	wantResponseCode int, wantResponseParsed T,
	gotResponse *httptest.ResponseRecorder) {

	assert.Equal(t, wantResponseCode, gotResponse.Code)

	if wantResponseCode >= 200 && wantResponseCode < 300 {
		var gotBodyParsed T
		err := json.NewDecoder(gotResponse.Body).Decode(&gotBodyParsed)
		if err != nil {
			t.Error("response is not valid json: %w", err)
			return
		}
		assert.Equal(t, wantResponseParsed, gotBodyParsed)
	} else if wantResponseCode == http.StatusInternalServerError {
		var gotBodyParsed errorResponses.ErrorResponse
		err := json.NewDecoder(gotResponse.Body).Decode(&gotBodyParsed)
		if err != nil {
			t.Error("response is not valid json: %w", err)
			return
		}
		assert.Equal(t, errorResponses.ErrorInternalResponse, gotBodyParsed)
	}
}

// NewJSONRequest wraps [http.NewRequest] and sets the passed struct as body
func NewJSONRequest(method, url string, bodyAsStruct any) (*http.Request, error) {
	body, err := json.Marshal(bodyAsStruct)
	if err != nil {
		return nil, fmt.Errorf("could not parse struct to json: %w", err)
	}
	return http.NewRequest(method, url, bytes.NewReader(body))
}

// MockAuthMiddleware mocks authentication by setting a default user context in the Gin context for testing purposes.
func MockAuthMiddleware(ctx *gin.Context) {
	const userContextKey = "USER_CONTEXT_DATA"
	const iamRoleUser = "user"

	userContext := auth.UserContext{
		Realm:          "",
		UserID:         "",
		UserName:       "",
		EmailAddress:   "",
		Roles:          []string{iamRoleUser},
		Groups:         nil,
		AllowedOrigins: nil,
	}

	ctx.Set(userContextKey, userContext)
	ctx.Next()
}
