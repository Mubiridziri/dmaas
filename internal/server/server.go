package server

import (
	"dmaas/docs"
	"dmaas/internal/usecase/dictionaries"
	"dmaas/internal/usecase/dictionarydata"
	"dmaas/internal/usecase/sources"
	"dmaas/internal/usecase/tabledata"
	"dmaas/internal/usecase/tables"
	"dmaas/internal/usecase/users"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

const UserKey = "AUTH"

// TODO move to env APP_SECRET
var secret = []byte("RHYaxoa6iqb1VTCsFtdM2PAAu8i8CYhU")

type Config struct {
	UserController           *users.Controller
	SourceController         *sources.Controller
	TableDataController      *tabledata.Controller
	TableController          *tables.Controller
	DictionaryController     *dictionaries.Controller
	DictionaryDataController *dictionarydata.Controller
}

type Server struct {
	Router                   *gin.Engine
	userController           *users.Controller
	sourceController         *sources.Controller
	tableDataController      *tabledata.Controller
	tableController          *tables.Controller
	dictionaryController     *dictionaries.Controller
	dictionaryDataController *dictionarydata.Controller
}

type ListQuery struct {
	Page     int
	Limit    int
	Simplify bool
}

func New(config Config) *Server {
	s := Server{
		Router:                   gin.New(),
		userController:           config.UserController,
		sourceController:         config.SourceController,
		tableDataController:      config.TableDataController,
		tableController:          config.TableController,
		dictionaryController:     config.DictionaryController,
		dictionaryDataController: config.DictionaryDataController,
	}

	s.registerRoutes()
	return &s
}

func (s *Server) registerRoutes() {
	//Middleware
	s.Router.Use(gin.Logger())
	s.Router.Use(gin.Recovery())
	s.Router.Use(sessions.Sessions(UserKey, cookie.NewStore(secret)))

	//Swagger
	configureSwagger(s.Router)

	// K8s probe
	s.Router.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	//API /api/v1
	mainGroup := s.Router.Group("/api/v1")
	mainGroup.POST("/login", s.handleLogin)

	mainGroup.Use(AuthRequired(s.userController))
	{
		s.AddUserRoutes(mainGroup)
		s.AddSourceRoutes(mainGroup)
		s.AddDictionaryRoutes(mainGroup)

		mainGroup.GET("/login", s.handleProfile)
		mainGroup.GET("/logout", s.handleLogout)
	}
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

func AuthRequired(controller *users.Controller) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userKey := session.Get(UserKey)

		if userKey == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		authUser, err := controller.GetUserByUsername(userKey.(string))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		c.Set("user", authUser)
		c.Next()
	}

}
