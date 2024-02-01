package v1

import (
	"dmaas/internal/context"
	"dmaas/internal/controller/response"
	"dmaas/internal/dto"
	"dmaas/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TableController struct {
	Context *context.ApplicationContext
}

type PaginatedTables struct {
	Total   int64          `json:"total"`
	Entries []entity.Table `json:"entries"`
}

type PaginatedTablesData struct {
	Total   int64                    `json:"total"`
	Entries []map[string]interface{} `json:"entries"`
}

// listTablesAction GoDoc
//
//	@Summary	List Table
//	@Schemes
//	@Description	Paginated Table List
//	@Param			sourceId	path	int	true	"SourceID"
//	@Param			page	query	int	false	"Page"
//	@Param			limit	query	int	false	"Limit"
//	@Tags			Sources
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	PaginatedTables
//	@Router			/api/v1/sources/:id/tables [GET]
func (controller *TableController) listTablesAction(c *gin.Context) {
	sourceId, err := dto.IdFromContext(c)

	if err != nil {
		response.CreateBadRequestResponse(c, err.Error())
	}

	pagination, err := dto.QueryFromContext(c)

	if err != nil {
		response.CreateBadRequestResponse(c, err.Error())
	}

	entries, err := controller.Context.TableUseCase.ListTables(sourceId, pagination)
	count := controller.Context.TableUseCase.GetCount()

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, PaginatedTables{
		Total:   count,
		Entries: entries,
	})
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
//	@Success		200	{array}	PaginatedTablesData
//	@Router			/api/v1/sources/table/data/:id [GET]
func (controller *TableController) listTableDataAction(c *gin.Context) {
	tableId, err := dto.IdFromContext(c)

	if err != nil {
		response.CreateBadRequestResponse(c, err.Error())
	}

	pagination, err := dto.QueryFromContext(c)

	if err != nil {
		response.CreateBadRequestResponse(c, err.Error())
	}

	table, err := controller.Context.TableUseCase.GetTableById(tableId)

	if err != nil {
		response.CreateNotFoundResponse(c, "not found")
		return
	}

	source, err := controller.Context.SourceUseCase.GetSourceById(table.SourceID)

	if err != nil {
		response.CreateNotFoundResponse(c, "not found")
		return
	}

	data, err := controller.Context.TableDataUseCase.ListTableData(source, table, pagination)
	count := controller.Context.TableDataUseCase.GetCount(source, table)

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, PaginatedTablesData{
		Total:   count,
		Entries: data,
	})
}

func (controller *TableController) AddTableRoute(r *gin.RouterGroup) {
	userGroup := r.Group("/sources")

	userGroup.GET("/:id/tables", controller.listTablesAction)
	userGroup.GET("/table/data/:id", controller.listTableDataAction)
}
