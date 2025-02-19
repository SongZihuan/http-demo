package engine

import (
	"github.com/SongZihuan/http-demo/src/handler"
	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine = nil

func InitEngine() error {
	gin.SetMode(gin.ReleaseMode)

	Engine = gin.New()
	Engine.Use(gin.Logger(), gin.Recovery())

	Engine.GET("/", handler.HandlerMessage)
	Engine.GET("/message", handler.HandlerMessage)
	Engine.GET("/ip", handler.HandlerRemoteIP)
	Engine.GET("/remote/ip", handler.HandlerRemoteIP)
	Engine.GET("/client/ip", handler.HandlerClientIP)
	Engine.GET("/timestamp", handler.HandlerTimestamp)
	Engine.GET("/datetime", handler.HandlerDatetime)
	Engine.GET("/hello", handler.HandlerHelloWorld)
	Engine.GET("/empty", handler.HandlerEmpty)
	Engine.NoRoute(handler.HandlerMethodNotFound)
	Engine.NoMethod(handler.HandlerMethodNotAllowed)

	return nil
}
