package routes

import (
	"dmaas/internal/app/dmaas/repository"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Repository repository.UserRepositoryInterface
}

func (controller *UserController) listUsersAction(c *gin.Context)  {}
func (controller *UserController) createUserAction(c *gin.Context) {}
func (controller *UserController) editUserAction(c *gin.Context)   {}
func (controller *UserController) detailUserAction(c *gin.Context) {}
func (controller *UserController) removeUserAction(c *gin.Context) {}

func (controller *UserController) AddUserRoute(r *gin.RouterGroup) {
	userGroup := r.Group("/users")

	userGroup.GET("", controller.listUsersAction)
	userGroup.POST("", controller.createUserAction)
	userGroup.PUT("/:id", controller.editUserAction)
	userGroup.GET("/:id", controller.detailUserAction)
	userGroup.DELETE("/:id", controller.removeUserAction)
}
