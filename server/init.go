package server

import "github.com/gin-gonic/gin"

var Server *gin.Engine

func init() {
	Server := gin.New()
	Server.Use(gin.Logger())
	Server.Use(gin.Recovery())
	userRouter := Server.Group("/api/user")
	scriptRouter := Server.Group("/api/script")
	loadUserRouter(userRouter)
	loadScriptRouter(scriptRouter)
}
