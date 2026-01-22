package testhelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/greenbone/opensight-notification-service/pkg/helper"
	"github.com/greenbone/opensight-notification-service/pkg/pgtesting"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/greenbone/opensight-notification-service/pkg/repository/notificationrepository"
	"github.com/greenbone/opensight-notification-service/pkg/request"
	"github.com/greenbone/opensight-notification-service/pkg/security"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func VerifyResponseWithMetadata[T any](
	t *testing.T,
	wantResponseCode int, wantResponseParsed T,
	gotResponse *httptest.ResponseRecorder,
) {

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

// MockAuthMiddleware mocks authentication by setting a admin user context in the Gin context for testing purposes.
func MockAuthMiddlewareWithAdmin(ctx *gin.Context) {
	const userContextKey = "USER_CONTEXT_DATA"
	const iamRoleUser = "admin"

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

func SetupNotificationChannelTestEnv(t *testing.T) (port.NotificationChannelRepository, *sqlx.DB) {
	encryptMgr := security.NewEncryptManager()
	encryptMgr.UpdateKeys(config.DatabaseEncryptionKey{
		Password:     "password",
		PasswordSalt: "password-salt-should-no-be-short-fyi",
	})

	db := pgtesting.NewDB(t)
	repo, err := notificationrepository.NewNotificationChannelRepository(db, encryptMgr)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	return repo, db
}

func GetValidMailNotificationChannel() request.MailNotificationChannelRequest {
	return request.MailNotificationChannelRequest{
		ChannelName:              "mail1",
		Domain:                   "example.com",
		Port:                     25,
		IsAuthenticationRequired: true,
		IsTlsEnforced:            false,
		Username:                 helper.ToPtr("user"),
		Password:                 helper.ToPtr("pass"),
		MaxEmailAttachmentSizeMb: helper.ToPtr(10),
		MaxEmailIncludeSizeMb:    helper.ToPtr(5),
		SenderEmailAddress:       "sender@example.com",
	}
}

func GetValidMattermostNotificationChannel() request.MattermostNotificationChannelRequest {
	return request.MattermostNotificationChannelRequest{
		ChannelName: "mattermost1",
		WebhookUrl:  "http://webhookurl.com/id1",
		Description: "This is a test mattermost channel",
	}
}
