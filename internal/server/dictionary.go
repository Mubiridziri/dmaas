package server

import (
	"dmaas/internal/usecase/dictionaries"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (s *Server) handleListDictionaries(c *gin.Context) {
	var query ListQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad query params: " + err.Error(),
		})
		return
	}

	rows, err := s.dictionaryController.ListDictionaries(query.Page, query.Limit)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Server error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, rows)
}

func (s *Server) handleCreateDictionary(c *gin.Context) {
	var createDict dictionaries.CreateOrUpdateDictionaryView

	if err := c.ShouldBindJSON(&createDict); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Bind error: " + err.Error(),
		})
		return
	}

	outputSources, err := s.dictionaryController.CreateDictionary(createDict)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Server error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, outputSources)
}

func (s *Server) handleDetailDictionary(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID param in path",
		})
		return
	}

	outputSource, err := s.dictionaryController.GetDictionaryById(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Not found",
		})
		return
	}

	c.JSON(http.StatusOK, outputSource)
}

func (s *Server) handleDeleteDictionary(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID param in path",
		})
		return
	}

	outputSource, err := s.dictionaryController.RemoveDictionary(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Not found",
		})
		return
	}

	c.JSON(http.StatusOK, outputSource)
}

func (s *Server) AddDictionaryRoutes(g *gin.RouterGroup) {
	grp := g.Group("/dictionaries")

	grp.GET("", s.handleListDictionaries)
	grp.POST("", s.handleCreateDictionary)
	grp.GET("/:id", s.handleDetailDictionary)
	grp.DELETE("/:id", s.handleDeleteDictionary)

}
