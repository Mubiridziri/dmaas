package routes

import (
	"dmaas/internal/app/dmaas/controller/response"
	"dmaas/internal/app/dmaas/entity"
	"dmaas/internal/app/dmaas/repository"
	"dmaas/internal/app/dmaas/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type SourceController struct {
	Repository    repository.SourceRepositoryInterface
	SourceManager sources.SourceManager
}

type SourceRequest struct {
	Title    string `json:"title" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Schema   string `json:"schema" binding:"required"`
}

// Validate TODO MOVE TO .... ???
func (model *SourceRequest) Validate() error {
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

func (controller *SourceController) listSourcesAction(c *gin.Context) {}

func (controller *SourceController) createSourceAction(c *gin.Context) {
	var request SourceRequest

	if err := c.BindJSON(&request); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	//TODO MOVE TO .... ????
	source := entity.Source{
		Title:    request.Title,
		Name:     request.Name,
		Type:     request.Type,
		Host:     request.Host,
		Port:     request.Port,
		Username: request.Username,
		Password: request.Password,
		Schema:   request.Schema,
	}
	err := controller.Repository.CreateSource(&source)

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	//TODO MOVE TO ... ???
	go controller.SourceManager.ImportDatabase(source) //side task for import database

	c.JSON(http.StatusCreated, source)

}

func (controller *SourceController) editSourceAction(c *gin.Context) {}

func (controller *SourceController) detailSourceAction(c *gin.Context) {}

func (controller *SourceController) removeSourceAction(c *gin.Context) {}

func (controller *SourceController) AddSourceRoute(r *gin.RouterGroup) {
	userGroup := r.Group("/sources")

	userGroup.GET("", controller.listSourcesAction)
	userGroup.POST("", controller.createSourceAction)
	userGroup.PUT("/:id", controller.editSourceAction)
	userGroup.GET("/:id", controller.detailSourceAction)
	userGroup.DELETE("/:id", controller.removeSourceAction)
}
