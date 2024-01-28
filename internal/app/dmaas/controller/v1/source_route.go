package v1

import (
	"dmaas/internal/app/dmaas/controller/response"
	"dmaas/internal/app/dmaas/entity"
	"dmaas/internal/app/dmaas/repository"
	"dmaas/internal/app/dmaas/service"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type SourceController struct {
	Repository    repository.SourceRepositoryInterface
	SourceManager *sources.SourceManager
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

type PaginatedSources struct {
	Total   int64           `json:"total"`
	Entries []entity.Source `json:"entries"`
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

// listSourcesAction GoDoc
//
//	@Summary	List Source
//	@Schemes
//	@Description	Paginated Source List
//	@Param			page	query	int	false	"Page"
//	@Param			limit	query	int	false	"Page"
//	@Tags			Sources
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	PaginatedSources
//	@Router			/api/v1/sources [GET]
func (controller *SourceController) listSourcesAction(c *gin.Context) {
	//TODO may be bind to model (struct) ?
	pageQuery := c.DefaultQuery("page", "1")
	limitQuery := c.DefaultQuery("limit", "10")

	page, pageOk := strconv.Atoi(pageQuery)
	limit, limitOk := strconv.Atoi(limitQuery)

	if pageOk != nil || limitOk != nil {
		response.CreateBadRequestResponse(c, "bad query parameters")
		return
	}

	entries, err := controller.Repository.ListSources(page, limit)
	count := controller.Repository.GetCount()

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, PaginatedSources{
		Total:   count,
		Entries: entries,
	})
}

// createSourceAction GoDoc
//
//	@Summary	Create Source
//	@Schemes
//	@Description	Create entity
//	@Param			request	body	SourceRequest	true	"Source Data"
//	@Tags			Sources
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Source
//	@Router			/api/v1/sources [POST]
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
	//boilerplate
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

// editSourceAction GoDoc
//
//	@Summary	Update Source
//	@Schemes
//	@Description	Update entity
//	@Tags			Sources
//	@Param			id		path	int				true	"Source ID"
//	@Param			request	body	SourceRequest	true	"Source Data"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Source
//	@Router			/api/v1/sources/:id [PUT]
func (controller *SourceController) editSourceAction(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		response.CreateNotFoundResponse(c, "invalid ID param")
		return
	}

	source, err := controller.Repository.GetSourceById(id)

	if err != nil {
		response.CreateNotFoundResponse(c, "invalid ID param")
		return
	}

	var request SourceRequest

	if err := c.BindJSON(&request); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	//TODO Refact?? need object to populate method
	//boilerplate
	source.Title = request.Title
	source.Name = request.Name
	source.Type = request.Type
	source.Host = request.Host
	source.Port = request.Port
	source.Username = request.Username
	source.Password = request.Password
	source.Schema = request.Schema

	err = controller.Repository.UpdateSource(&source)

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, source)

}

// detailSourceAction GoDoc
//
//	@Summary	Detail Source
//	@Schemes
//	@Description	Get By ID
//	@Tags			Sources
//	@Param			id	path	int	true	"Source ID"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Source
//	@Router			/api/v1/sources/:id [GET]
func (controller *SourceController) detailSourceAction(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		response.CreateNotFoundResponse(c, "invalid ID param")
		return
	}

	source, err := controller.Repository.GetSourceById(id)

	if err != nil {
		response.CreateNotFoundResponse(c, "not found")
		return
	}

	c.JSON(http.StatusOK, source)
}

// removeSourceAction GoDoc
//
//	@Summary	Remove Source
//	@Schemes
//	@Description	Remove By ID
//	@Tags			Sources
//	@Param			id	path	int	true	"Source ID"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Source
//	@Router			/api/v1/sources/:id [DELETE]
func (controller *SourceController) removeSourceAction(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		response.CreateNotFoundResponse(c, "invalid ID param")
		return
	}

	source, err := controller.Repository.GetSourceById(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.CreateNotFoundResponse(c, "not found")
		return
	}

	err = controller.Repository.RemoveSource(&source)

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	go controller.SourceManager.DeleteDatabase(source)

	c.JSON(http.StatusOK, source)
}

func (controller *SourceController) AddSourceRoute(r *gin.RouterGroup) {
	userGroup := r.Group("/sources")

	userGroup.GET("", controller.listSourcesAction)
	userGroup.POST("", controller.createSourceAction)
	userGroup.PUT("/:id", controller.editSourceAction)
	userGroup.GET("/:id", controller.detailSourceAction)
	userGroup.DELETE("/:id", controller.removeSourceAction)
}
