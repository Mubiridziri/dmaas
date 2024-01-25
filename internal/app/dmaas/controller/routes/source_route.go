package routes

import (
	"dmaas/internal/app/dmaas/controller/response"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

type NewSourceRequest struct {
	Title    string `json:"title" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Schema   string `json:"schema" binding:"required"`
}

// Validate Where need move validation for models?
func (model *NewSourceRequest) Validate() error {
	if strings.Contains(model.Name, ";") {
		return errors.New("name contains is forbidden symbol")
	}

	if strings.Contains(model.Host, ";") {
		return errors.New("host contains is forbidden symbol")
	}
	if strings.Contains(model.Username, ";") {
		return errors.New("username contains is forbidden symbol")
	}
	if strings.Contains(model.Password, ";") {
		return errors.New("password contains is forbidden symbol")
	}
	if strings.Contains(model.Schema, ";") {
		return errors.New("schema contains is forbidden symbol")
	}

	return nil
}

func listSourcesAction(c *gin.Context) {}

func createSourceAction(c *gin.Context) {
	var request NewSourceRequest

	if err := c.BindJSON(&request); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

}

func editSourceAction(c *gin.Context) {}

func detailSourceAction(c *gin.Context) {}

func removeSourceAction(c *gin.Context) {}

func AddSourceRoute(r *gin.RouterGroup) {
	userGroup := r.Group("/sources")

	userGroup.GET("", listSourcesAction)
	userGroup.POST("", createSourceAction)
	userGroup.PUT("/:id", editSourceAction)
	userGroup.GET("/:id", detailSourceAction)
	userGroup.DELETE("/:id", removeSourceAction)
}
