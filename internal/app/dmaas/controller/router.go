package controller

import (
	"dmaas/internal/app/dmaas/controller/middleware"
	"dmaas/internal/app/dmaas/controller/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// TODO generate when Application Run
var secret = []byte("RHYaxoa6iqb1VTCsFtdM2PAAu8i8CYhU")

func NewRouter() *gin.Engine {
	r := gin.New()

	//Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(sessions.Sessions("AUTH", cookie.NewStore(secret)))

	//API
	v1 := r.Group("/api/v1")
	v1.POST("/login", routes.LoginAction)

	v1.Use(middleware.AuthRequired)
	{
		v1.GET("/logout", routes.LogoutAction)
		routes.AddUserRoute(v1)
		routes.AddSourceRoute(v1)
	}

	return r
}
