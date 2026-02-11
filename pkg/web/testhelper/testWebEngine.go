package testhelper

import (
	"github.com/gin-gonic/gin"
	logsMiddleware "github.com/greenbone/opensight-golang-libraries/pkg/logs/ginMiddleware"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

func NewTestWebEngine(registry *errmap.Registry) *gin.Engine {
	gin.SetMode(gin.TestMode)

	ginWebEngine := gin.New()
	ginWebEngine.Use(
		logsMiddleware.Logging(),
		gin.Recovery(),
		middleware.CORS([]string{
			"http://example.com",
		}),
		middleware.ErrorHandler(gin.ErrorTypeAny),
		middleware.InterpretErrors(gin.ErrorTypePrivate, registry),
	)

	return ginWebEngine
}
