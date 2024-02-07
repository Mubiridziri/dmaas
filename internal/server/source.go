package server

import (
	"dmaas/internal/usecase/sources"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// handleListSource GoDoc
//
//	@Summary	Get Sources List
//	@Schemes
//	@Description	List of sources
//	@Param			page	query	int	false	"Page"
//	@Param			limit	query	int	false	"Limit of page"
//	@Tags			Sources
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]sources.PaginatedSourcesList
//	@Router			/api/v1/sources [GET]
func (s *Server) handleListSource(c *gin.Context) {
	var query ListQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad query params: " + err.Error(),
		})
		return
	}

	rows, err := s.sourceController.ListSources(query.Page, query.Limit)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Server error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, rows)
}

// handleCreateSource GoDoc
//
//	@Summary	Create Source
//	@Schemes
//	@Description	Creating source
//	@Param			source	body	sources.CreateOrUpdateSourceView	true	"Source"
//	@Tags			Sources
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	sources.SourceView
//	@Router			/api/v1/sources [POST]
func (s *Server) handleCreateSource(c *gin.Context) {
	var createSource sources.CreateOrUpdateSourceView

	if err := c.ShouldBindJSON(&createSource); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Bind error: " + err.Error(),
		})
		return
	}

	outputSources, err := s.sourceController.CreateSource(createSource)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Server error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, outputSources)
}

// handleDetailSource GoDoc
//
//	@Summary	Detail Source
//	@Schemes
//	@Description	Get source info by source id
//	@Param			id	path	int	true	"Source ID"
//	@Tags			Sources
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	sources.SourceView
//	@Router			/api/v1/sources/{id} [GET]
func (s *Server) handleDetailSource(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID param in path",
		})
		return
	}

	outputSource, err := s.sourceController.GetSourceById(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Not found",
		})
		return
	}

	c.JSON(http.StatusOK, outputSource)
}

// handleDeleteSource GoDoc
//
//	@Summary	Delete Source
//	@Schemes
//	@Description	Deleting source
//	@Param			id	path	int	true	"Source ID"
//	@Tags			Sources
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	sources.SourceView
//	@Router			/api/v1/sources/{id} [DELETE]
func (s *Server) handleDeleteSource(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID param in path",
		})
		return
	}

	outputSource, err := s.sourceController.RemoveSource(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Not found",
		})
		return
	}

	c.JSON(http.StatusOK, outputSource)
}

func (s *Server) AddSourceRoutes(g *gin.RouterGroup) {
	grp := g.Group("/sources")

	grp.GET("", s.handleListSource)
	grp.POST("", s.handleCreateSource)
	grp.GET("/:id", s.handleDetailSource)
	grp.DELETE("/:id", s.handleDeleteSource)
}
