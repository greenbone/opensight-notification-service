package testhelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/golangmigrator"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

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

// NewTestDB creates a new database for integration testing
func NewTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	gm := golangmigrator.New(getMigrationsPath())
	cfg := pgtestdb.Config{
		DriverName: "pgx",
		User:       "postgres",
		Password:   "password",
		Database:   "notification_service",
		TestRole: &pgtestdb.Role{
			Username: "postgres",
			Password: "password",
		},
		Host:                      "localhost",
		Port:                      "5432",
		Options:                   "sslmode=disable",
		ForceTerminateConnections: true,
	}

	sqlDb := pgtestdb.Custom(t, cfg, gm)

	url := sqlDb.URL() + "&search_path=application"
	pgConfig := postgres.Config{
		DriverName: cfg.DriverName,
		DSN:        url,
	}
	gormConfig := &gorm.Config{
		TranslateError: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	gormDB, err := gorm.Open(postgres.New(pgConfig), gormConfig)

	require.NoError(t, err)

	return gormDB
}

func getMigrationsPath() string {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	return filepath.Join(basePath, "../../repository/migrations")
}
