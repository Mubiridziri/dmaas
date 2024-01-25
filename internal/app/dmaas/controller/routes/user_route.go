package routes

import "github.com/gin-gonic/gin"

func listUsersAction(c *gin.Context)  {}
func createUserAction(c *gin.Context) {}
func editUserAction(c *gin.Context)   {}
func detailUserAction(c *gin.Context) {}
func removeUserAction(c *gin.Context) {}

func AddUserRoute(r *gin.RouterGroup) {
	userGroup := r.Group("/users")

	userGroup.GET("", listUsersAction)
	userGroup.POST("", createUserAction)
	userGroup.PUT("/:id", editUserAction)
	userGroup.GET("/:id", detailUserAction)
	userGroup.DELETE("/:id", removeUserAction)
}
