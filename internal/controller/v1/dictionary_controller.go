package v1

import (
	"dmaas/internal/context"
	"dmaas/internal/controller/response"
	"dmaas/internal/dto"
	"dmaas/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DictionaryController struct {
	Context *context.ApplicationContext
}

type PaginatedDictionaries struct {
	Total   int64               `json:"total"`
	Entries []entity.Dictionary `json:"entries"`
}

// listDictionariesAction GoDoc
//
//	@Summary	List Dictionary
//	@Schemes
//	@Description	Paginated Dictionary List
//	@Param			page	query	int	false	"Page"
//	@Param			limit	query	int	false	"Limit"
//	@Tags			Dictionaries
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	PaginatedDictionaries
//	@Router			/api/v1/dictionaries [GET]
func (controller *DictionaryController) listDictionariesAction(c *gin.Context) {
	pagination, err := dto.QueryFromContext(c)

	if err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	entries, err := controller.Context.DictionaryUseCase.ListDictionaries(pagination)
	count := controller.Context.DictionaryUseCase.GetCount()

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, PaginatedDictionaries{
		Total:   count,
		Entries: entries,
	})
}

// createDictionaryAction GoDoc
//
//	@Summary	Create Dictionary
//	@Schemes
//	@Description	Create entity
//	@Param			request	body	dto.DictionaryRequest	true	"Dictionary Data"
//	@Tags			Dictionaries
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Dictionary
//	@Router			/api/v1/dictionaries [POST]
func (controller *DictionaryController) createDictionaryAction(c *gin.Context) {
	var request dto.DictionaryRequest

	if err := c.BindJSON(&request); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	dictionary := request.ToDictionary()
	err := controller.Context.DictionaryUseCase.CreateDictionary(&dictionary)
	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	err = controller.Context.DictionaryUseCase.CreateDictionary(&dictionary)
	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, dictionary)
}

// editDictionaryAction GoDoc
//
//	@Summary	Update Dictionary
//	@Schemes
//	@Description	Update entity
//	@Tags			Dictionaries
//	@Param			id		path	int				true	"Dictionary ID"
//	@Param			request	body	dto.DictionaryUpdateRequest	true	"Dictionary Data"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Dictionary
//	@Router			/api/v1/dictionaries/:id [PUT]
func (controller *DictionaryController) editDictionaryAction(c *gin.Context) {
	id, err := dto.IdFromContext(c)

	if err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	dictionary, err := controller.Context.DictionaryUseCase.GetDictionaryById(id)

	if err != nil {
		response.CreateNotFoundResponse(c, "invalid ID param")
		return
	}

	var request dto.DictionaryUpdateRequest

	if err := c.BindJSON(&request); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	request.ToDictionary(&dictionary)

	//TODO нужно отслеживать те, что удалили!!
	err = controller.Context.DictionaryUseCase.UpdateDictionary(&dictionary)

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, dictionary)
}

// detailDictionaryAction GoDoc
//
//	@Summary	Detail Dictionary
//	@Schemes
//	@Description	Get By ID
//	@Tags			Dictionaries
//	@Param			id	path	int	true	"Dictionary ID"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Dictionary
//	@Router			/api/v1/dictionaries/:id [GET]
func (controller *DictionaryController) detailDictionaryAction(c *gin.Context) {
	id, err := dto.IdFromContext(c)

	if err != nil {
		response.CreateBadRequestResponse(c, "invalid ID param")
		return
	}

	dictionary, err := controller.Context.DictionaryUseCase.GetDictionaryById(id)

	if err != nil {
		response.CreateNotFoundResponse(c, "not found")
		return
	}

	c.JSON(http.StatusOK, dictionary)
}

// removeDictionaryAction GoDoc
//
//	@Summary	Remove Dictionary
//	@Schemes
//	@Description	Remove By ID
//	@Tags			Dictionaries
//	@Param			id	path	int	true	"Dictionary ID"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Dictionary
//	@Router			/api/v1/dictionaries/:id [DELETE]
func (controller *DictionaryController) removeDictionaryAction(c *gin.Context) {
	id, err := dto.IdFromContext(c)

	if err != nil {
		response.CreateBadRequestResponse(c, "invalid ID param")
		return
	}

	dictionary, err := controller.Context.DictionaryUseCase.GetDictionaryById(id)
	err = controller.Context.DictionaryUseCase.RemoveDictionary(&dictionary)

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, dictionary)
}

func (controller *DictionaryController) AddDictionaryRoute(r *gin.RouterGroup) {
	dictionaryGroup := r.Group("/dictionaries")

	dictionaryGroup.GET("", controller.listDictionariesAction)
	dictionaryGroup.POST("", controller.createDictionaryAction)
	dictionaryGroup.PUT("/:id", controller.editDictionaryAction)
	dictionaryGroup.GET("/:id", controller.detailDictionaryAction)
	dictionaryGroup.DELETE("/:id", controller.removeDictionaryAction)
}
