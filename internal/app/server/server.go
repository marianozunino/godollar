package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type GinHandler struct {
	Gin *gin.Engine
}

var Module = fx.Options(
	fx.Provide(registerGinHandler),
)

func registerGinHandler() *GinHandler {
	instance := gin.Default()

	handler := GinHandler{
		Gin: instance,
	}

	return &handler
}
