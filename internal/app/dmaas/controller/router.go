package controller

import (
	"dmaas/docs"
	"dmaas/internal/app/dmaas/controller/middleware"
	"dmaas/internal/app/dmaas/controller/v1"
	"dmaas/internal/app/dmaas/repository"
	sources "dmaas/internal/app/dmaas/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

// TODO generate when Application Run
var secret = []byte("RHYaxoa6iqb1VTCsFtdM2PAAu8i8CYhU")

type Router struct {
	//Repositories
	SourceRepository    repository.SourceRepositoryInterface
	UserRepository      repository.UserRepositoryInterface
	TableRepository     repository.TableRepositoryInterface
	TableDataRepository repository.TableDataRepositoryInterface

	//Services
	SourceManager *sources.SourceManager
}

func (router *Router) NewRouter() *gin.Engine {
	r := gin.New()

	//Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(sessions.Sessions("AUTH", cookie.NewStore(secret)))

	//Swagger
	configureSwagger(r)

	//Init controllers
	securityController := v1.SecurityController{Repository: router.UserRepository}
	sourceController := v1.SourceController{Repository: router.SourceRepository, SourceManager: router.SourceManager}
	userController := v1.UserController{Repository: router.UserRepository}
	tableController := v1.TableController{
		TableRepository:     router.TableRepository,
		TableDataRepository: router.TableDataRepository,
		SourceRepository:    router.SourceRepository,
		SourceManager:       router.SourceManager,
	}

	// K8s probe
	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	//API /api/v1
	mainGroup := r.Group("/api/v1")
	mainGroup.POST("/login", securityController.LoginAction)

	mainGroup.Use(middleware.AuthRequired(router.UserRepository))
	{
		mainGroup.GET("/login", securityController.ProfileAction)
		mainGroup.GET("/logout", securityController.LogoutAction)

		userController.AddUserRoute(mainGroup)
		sourceController.AddSourceRoute(mainGroup)
		tableController.AddTableRoute(mainGroup)
	}

	return r
}

func configureSwagger(r *gin.Engine) {
	//Swagger
	docs.SwaggerInfo.Title = "DMAAS"
	docs.SwaggerInfo.Description = "Data management and analytic system"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
