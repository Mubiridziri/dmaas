package controller

import (
	"dmaas/internal/app/dmaas/controller/middleware"
	"dmaas/internal/app/dmaas/controller/routes"
	"dmaas/internal/app/dmaas/repository"
	sources "dmaas/internal/app/dmaas/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// TODO generate when Application Run
var secret = []byte("RHYaxoa6iqb1VTCsFtdM2PAAu8i8CYhU")

type Router struct {
	//Repositories
	SourceRepository repository.SourceRepositoryInterface
	UserRepository   repository.UserRepositoryInterface

	//Services
	SourceManager *sources.SourceManager
}

func (router *Router) NewRouter() *gin.Engine {
	r := gin.New()

	//Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(sessions.Sessions("AUTH", cookie.NewStore(secret)))

	//Init controllers
	securityController := routes.SecurityController{Repository: router.UserRepository}
	sourceController := routes.SourceController{Repository: router.SourceRepository}
	userController := routes.UserController{Repository: router.UserRepository}

	//API
	v1 := r.Group("/api/v1")
	v1.POST("/login", securityController.LoginAction)

	v1.Use(middleware.AuthRequired(router.UserRepository))
	{
		v1.GET("/logout", securityController.LogoutAction)

		userController.AddUserRoute(v1)
		sourceController.AddSourceRoute(v1)
	}

	return r
}
