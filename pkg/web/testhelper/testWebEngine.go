package testhelper

import (
	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/web/helper"
)

func NewTestWebEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(helper.ValidationErrorHandler(gin.ErrorTypePrivate))

	return router
}
