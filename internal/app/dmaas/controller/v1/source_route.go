package v1

import (
	"dmaas/internal/app/dmaas/context"
	"dmaas/internal/app/dmaas/controller/response"
	"dmaas/internal/app/dmaas/dto"
	"dmaas/internal/app/dmaas/entity"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type SourceController struct {
	Context *context.ApplicationContext
}

type PaginatedSources struct {
	Total   int64           `json:"total"`
	Entries []entity.Source `json:"entries"`
}

// listSourcesAction GoDoc
//
//	@Summary	List Source
//	@Schemes
//	@Description	Paginated Source List
//	@Param			page	query	int	false	"Page"
//	@Param			limit	query	int	false	"Limit"
//	@Tags			Sources
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	PaginatedSources
//	@Router			/api/v1/sources [GET]
func (controller *SourceController) listSourcesAction(c *gin.Context) {
	pagination, err := dto.QueryFromContext(c)

	if err != nil {
		response.CreateBadRequestResponse(c, err.Error())
	}

	entries, err := controller.Context.SourceUseCase.ListSources(pagination)
	count := controller.Context.SourceUseCase.GetCount()

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
//	@Param			request	body	dto.SourceRequest	true	"Source Data"
//	@Tags			Sources
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Source
//	@Router			/api/v1/sources [POST]
func (controller *SourceController) createSourceAction(c *gin.Context) {
	var request dto.SourceRequest

	if err := c.BindJSON(&request); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	source := request.ToSource()
	err := controller.Context.SourceUseCase.CreateSource(&source)

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, source)

}

// editSourceAction GoDoc
//
//	@Summary	Update Source
//	@Schemes
//	@Description	Update entity
//	@Tags			Sources
//	@Param			id		path	int				true	"Source ID"
//	@Param			request	body	dto.SourceRequest	true	"Source Data"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Source
//	@Router			/api/v1/sources/:id [PUT]
func (controller *SourceController) editSourceAction(c *gin.Context) {
	id, err := dto.IdFromContext(c)

	if err != nil {
		response.CreateBadRequestResponse(c, "invalid ID param")
		return
	}

	source, err := controller.Context.SourceUseCase.GetSourceById(id)

	if err != nil {
		response.CreateNotFoundResponse(c, "not found")
		return
	}

	var request dto.SourceUpdateRequest

	if err := c.BindJSON(&request); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	request.ToSource(&source)
	err = controller.Context.SourceUseCase.UpdateSource(&source)

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
	id, err := dto.IdFromContext(c)

	if err != nil {
		response.CreateBadRequestResponse(c, "invalid ID param")
		return
	}

	source, err := controller.Context.SourceUseCase.GetSourceById(id)

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
	id, err := dto.IdFromContext(c)

	if err != nil {
		response.CreateBadRequestResponse(c, "invalid ID param")
		return
	}

	source, err := controller.Context.SourceUseCase.GetSourceById(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.CreateNotFoundResponse(c, "not found")
		return
	}

	err = controller.Context.SourceUseCase.RemoveSource(&source)

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

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
