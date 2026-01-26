package mailcontroller

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/mapper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port/mocks"
	mailmocks "github.com/greenbone/opensight-notification-service/pkg/port/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/request"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getValidNotificationChannel() models.NotificationChannel {
	name := "mail1"
	domain := "example.com"
	port := 25
	auth := true
	tls := true
	username := "user"
	password := "pass"
	maxAttach := 10
	maxInclude := 5
	sender := "sender@example.com"
	return models.NotificationChannel{
		ChannelName:              &name,
		Domain:                   &domain,
		Port:                     &port,
		IsAuthenticationRequired: &auth,
		IsTlsEnforced:            &tls,
		Username:                 &username,
		Password:                 &password,
		MaxEmailAttachmentSizeMb: &maxAttach,
		MaxEmailIncludeSizeMb:    &maxInclude,
		SenderEmailAddress:       &sender,
	}
}

func setupRouter(service *mocks.NotificationChannelService, mailService *mailmocks.MailChannelService) *gin.Engine {
	registry := errmap.NewRegistry()

	engine := testhelper.NewTestWebEngine()
	engine.Use(middleware.InterpretErrors(gin.ErrorTypePrivate, registry))

	NewMailController(engine, service, mailService, testhelper.MockAuthMiddlewareWithAdmin, registry)
	return engine
}

func TestMailController_CreateMailChannel(t *testing.T) {
	valid := getValidNotificationChannel()
	mailValid := mapper.MapNotificationChannelToMail(valid)
	created := mailValid // Simulate returned object

	tests := []struct {
		name           string
		input          any
		mockReturn     request.MailNotificationChannelRequest
		mockErr        error
		wantStatusCode int
	}{
		{
			name:           "success",
			input:          mailValid,
			mockReturn:     created,
			wantStatusCode: http.StatusCreated,
		},
		{
			name:           "empty request",
			input:          `{}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "invalid input (missing required)",
			input:          struct{ Foo string }{"bar"},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "internal error",
			input:          mailValid,
			mockErr:        assert.AnError,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "empty body",
			input:          nil,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid sender email",
			input: func() request.MailNotificationChannelRequest {
				invalid := mailValid
				invalid.SenderEmailAddress = "not-an-email"
				return invalid
			}(),
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewNotificationChannelService(t)
			mockMailService := mailmocks.NewMailChannelService(t)
			router := setupRouter(mockService, mockMailService)

			if tt.wantStatusCode == http.StatusCreated || tt.wantStatusCode == http.StatusInternalServerError {
				mockMailService.On("CreateMailChannel", mock.Anything, mock.Anything).
					Return(tt.mockReturn, tt.mockErr).
					Once()
			}

			req := httpassert.New(t, router).Post("/notification-channel/mail")
			if tt.input != nil {
				req.JsonContentObject(tt.input)
			}
			resp := req.Expect()
			resp.StatusCode(tt.wantStatusCode)
			if tt.wantStatusCode == http.StatusCreated {
				resp.JsonPath("$.channelName", mailValid.ChannelName)
			}
			if tt.name == "internal error" {
				resp.JsonPath("$.title", httpassert.Matcher(func(
					t *testing.T,
					actual any,
				) bool {
					return actual != ""
				}))
			} else if tt.wantStatusCode == http.StatusBadRequest {
				if tt.name == "invalid sender email" {
					resp.JsonPath("$", httpassert.Matcher(func(t *testing.T, actual any) bool {
						m, ok := actual.(map[string]interface{})
						if !ok {
							return false
						}
						_, exists := m["senderEmailAddress"]
						return exists
					}))
				} else {
					resp.JsonPath("$", httpassert.Matcher(func(t *testing.T, actual any) bool { return actual != nil }))
				}
			}
			mockMailService.AssertExpectations(t)
		})
	}
}

func TestMailController_ListMailChannelsByType(t *testing.T) {
	valid := getValidNotificationChannel()
	channels := []models.NotificationChannel{valid}

	tests := []struct {
		name           string
		queryType      models.ChannelType
		mockReturn     []models.NotificationChannel
		mockErr        error
		wantStatusCode int
	}{
		{
			name:           "list by type",
			queryType:      models.ChannelTypeMail,
			mockReturn:     channels,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "internal error",
			queryType:      models.ChannelTypeMail,
			mockErr:        errors.New("db error"),
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "empty result",
			queryType:      models.ChannelTypeMail,
			mockReturn:     []models.NotificationChannel{},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewNotificationChannelService(t)
			router := setupRouter(mockService, nil)

			mockService.On("ListNotificationChannelsByType", mock.Anything, tt.queryType).
				Return(tt.mockReturn, tt.mockErr).
				Once()

			req := httpassert.New(t, router).Get("/notification-channel/mail")
			resp := req.Expect()
			resp.StatusCode(tt.wantStatusCode)
			if tt.wantStatusCode == http.StatusOK {
				resp.JsonPath("$", httpassert.HasSize(len(tt.mockReturn)))
			}
			if tt.wantStatusCode == http.StatusInternalServerError {
				resp.JsonPath("$.title", httpassert.Matcher(func(
					t *testing.T,
					actual any,
				) bool {
					return actual != ""
				}))
				resp.JsonPath("$", httpassert.HasSize(2))
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestMailController_UpdateMailChannel(t *testing.T) {
	valid := getValidNotificationChannel()
	updated := valid
	id := "mail-id-1"

	tests := []struct {
		name           string
		id             string
		input          any
		mockReturn     models.NotificationChannel
		mockErr        error
		wantStatusCode int
	}{
		{
			name:           "success",
			id:             id,
			input:          valid,
			mockReturn:     updated,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "invalid input (missing required)",
			id:             id,
			input:          struct{ Foo string }{"bar"},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "internal error",
			id:             id,
			input:          valid,
			mockErr:        errors.New("db error"),
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "empty body",
			id:             id,
			input:          nil,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewNotificationChannelService(t)
			mockMailService := mailmocks.NewMailChannelService(t)
			router := setupRouter(mockService, mockMailService)

			if tt.wantStatusCode == http.StatusOK || tt.wantStatusCode == http.StatusInternalServerError {
				mockService.On("UpdateNotificationChannel", mock.Anything, tt.id, mock.Anything).
					Return(tt.mockReturn, tt.mockErr).
					Once()
			}
			req := httpassert.New(t, router).Put("/notification-channel/mail/" + tt.id)
			if tt.input != nil {
				req.JsonContentObject(tt.input)
			}
			resp := req.Expect()
			resp.StatusCode(tt.wantStatusCode)
			if tt.wantStatusCode == http.StatusOK {
				resp.JsonPath("$.channelName", *valid.ChannelName)
			}
			if tt.wantStatusCode == http.StatusBadRequest {
				resp.JsonPath("$", httpassert.Matcher(func(t *testing.T, actual any) bool { return actual != nil }))
			}
			if tt.wantStatusCode == http.StatusInternalServerError {
				resp.JsonPath("$.title", httpassert.Matcher(func(
					t *testing.T,
					actual any,
				) bool {
					return actual != ""
				}))
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestMailController_DeleteMailChannel(t *testing.T) {
	id := "mail-id-1"

	tests := []struct {
		name           string
		id             string
		mockErr        error
		wantStatusCode int
	}{
		{
			name:           "success",
			id:             id,
			wantStatusCode: http.StatusNoContent,
		},
		{
			name:           "internal error",
			id:             id,
			mockErr:        errors.New("db error"),
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "non-existent id",
			id:             "notfound",
			mockErr:        errors.New("not found"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewNotificationChannelService(t)
			router := setupRouter(mockService, nil)

			mockService.On("DeleteNotificationChannel", mock.Anything, tt.id).
				Return(tt.mockErr).
				Once()

			req := httpassert.New(t, router).Delete("/notification-channel/mail/" + tt.id)
			resp := req.Expect()
			resp.StatusCode(tt.wantStatusCode)
			if tt.wantStatusCode == http.StatusNoContent {
				resp.NoContent()
			}
			if tt.wantStatusCode == http.StatusInternalServerError {
				resp.JsonPath("$.title", httpassert.Matcher(func(
					t *testing.T,
					actual any,
				) bool {
					return actual != ""
				}))
				resp.JsonPath("$", httpassert.HasSize(2))
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestValidateEmailAddress(t *testing.T) {
	channel := request.MailNotificationChannelRequest{
		Domain:      "example.com",
		Port:        587,
		ChannelName: "Test Channel",
	}

	t.Run("valid email", func(t *testing.T) {
		channel.SenderEmailAddress = "test@example.com"
		err := channel.Validate()
		assert.Empty(t, err["senderEmailAddress"])
	})

	t.Run("empty email", func(t *testing.T) {
		channel.SenderEmailAddress = ""
		err := channel.Validate()
		assert.Contains(t, err["senderEmailAddress"], "required")
	})

	t.Run("invalid email", func(t *testing.T) {
		channel.SenderEmailAddress = "invalid-email"
		err := channel.Validate()
		assert.Error(t, err)
		assert.Contains(t, err["senderEmailAddress"], "invalid")
	})
}

func TestValidateFields(t *testing.T) {
	validChannel := request.MailNotificationChannelRequest{
		Domain:             "example.com",
		Port:               587,
		SenderEmailAddress: "sender@example.com",
		ChannelName:        "Test Channel",
	}

	t.Run("all valid fields", func(t *testing.T) {
		errMap := validChannel.Validate()
		assert.Nil(t, errMap)
	})

	t.Run("missing domain", func(t *testing.T) {
		ch := validChannel
		ch.Domain = ""
		err := ch.Validate()
		assert.Contains(t, err["domain"], "required")
	})

	t.Run("invalid port", func(t *testing.T) {
		ch := validChannel
		ch.Port = 0
		err := ch.Validate()
		assert.Contains(t, err["port"], "required")
	})

	t.Run("missing channelName", func(t *testing.T) {
		ch := validChannel
		ch.ChannelName = ""
		err := ch.Validate()
		assert.Contains(t, err["channelName"], "required")
	})

	t.Run("multiple errors", func(t *testing.T) {
		invalidChannel := request.MailNotificationChannelRequest{}
		err := invalidChannel.Validate()
		assert.NotNil(t, err)
		assert.Contains(t, err["domain"], "required")
		assert.Contains(t, err["port"], "required")
		assert.Contains(t, err["senderEmailAddress"], "required")
		assert.Contains(t, err["channelName"], "required")
	})
}
