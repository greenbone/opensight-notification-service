package testhelper

import (
	"github.com/gin-gonic/gin"
)

func NewTestWebEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}
