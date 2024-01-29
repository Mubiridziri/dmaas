package dto

import (
	"dmaas/internal/app/dmaas/controller/response"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Query struct {
	Page     int
	Limit    int
	Simplify bool //for simple response for lists (ex: selects)
}

func QueryFromContext(c *gin.Context) (Query, error) {
	var query Query

	if err := c.BindQuery(&query); err != nil {
		return Query{}, errors.New("invalid query parameters")
	}

	return query, nil
}

func IdFromContext(c *gin.Context) (int, error) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		response.CreateNotFoundResponse(c, "")
		return 0, errors.New("invalid ID param")
	}

	return id, nil
}
