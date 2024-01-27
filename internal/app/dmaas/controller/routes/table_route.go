package routes

import (
	"dmaas/internal/app/dmaas/controller/response"
	"dmaas/internal/app/dmaas/repository"
	sources "dmaas/internal/app/dmaas/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TableController struct {
	TableRepository     repository.TableRepositoryInterface
	TableDataRepository repository.TableDataRepositoryInterface
	SourceRepository    repository.SourceRepositoryInterface
	SourceManager       *sources.SourceManager
}

// listTablesAction GoDoc
//
//	@Summary	List Table
//	@Schemes
//	@Description	Paginated Table List
//	@Param			sourceId	path	int	true	"SourceID"
//	@Param			page	query	int	false	"Page"
//	@Param			limit	query	int	false	"Page"
//	@Tags			Sources
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	entity.Table
//	@Router			/api/v1/sources/:id/tables [GET]
func (controller *TableController) listTablesAction(c *gin.Context) {
	idParam := c.Param("id")

	sourceId, err := strconv.Atoi(idParam)

	//TODO may be bind to model (struct) ?
	pageQuery := c.DefaultQuery("page", "1")
	limitQuery := c.DefaultQuery("page", "10")

	page, pageOk := strconv.Atoi(pageQuery)
	limit, limitOk := strconv.Atoi(limitQuery)

	if pageOk != nil || limitOk != nil {
		response.CreateBadRequestResponse(c, "bad query parameters")
		return
	}

	entries, err := controller.TableRepository.ListTables(sourceId, page, limit)

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, entries)
}

// listTablesAction GoDoc
//
//	@Summary	List Table Data
//	@Schemes
//	@Description	Paginated Table List
//	@Param			sourceId	path	int	true	"SourceID"
//	@Param			sourceId	path	int	true	"TableID"
//	@Param			page	query	int	false	"Page"
//	@Param			limit	query	int	false	"Limit"
//	@Tags			Sources
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	entity.Table
//	@Router			/api/v1/sources/table/data/:id [GET]
func (controller *TableController) listTableDataAction(c *gin.Context) {
	idParam := c.Param("id")

	tableId, err := strconv.Atoi(idParam)

	//TODO may be bind to model (struct) ?
	pageQuery := c.DefaultQuery("page", "1")
	limitQuery := c.DefaultQuery("page", "10")

	page, pageOk := strconv.Atoi(pageQuery)
	limit, limitOk := strconv.Atoi(limitQuery)

	if pageOk != nil || limitOk != nil {
		response.CreateBadRequestResponse(c, "bad query parameters")
		return
	}

	table, err := controller.TableRepository.GetTableById(tableId)

	if err != nil {
		response.CreateNotFoundResponse(c, "not found")
		return
	}

	source, err := controller.SourceRepository.GetSourceById(table.SourceID)

	if err != nil {
		response.CreateNotFoundResponse(c, "not found")
		return
	}

	localSchemaName := controller.SourceManager.GetLocalSchemaName(source)

	data, err := controller.TableDataRepository.ListTableData(localSchemaName, table, page, limit)

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

func (controller *TableController) AddTableRoute(r *gin.RouterGroup) {
	userGroup := r.Group("/sources")

	userGroup.GET("/:id/tables", controller.listTablesAction)
	userGroup.GET("/table/data/:id", controller.listTableDataAction)
}
