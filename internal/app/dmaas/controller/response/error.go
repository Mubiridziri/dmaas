package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func CreateBadRequestResponse(c *gin.Context, text string) {
	createResponse(c, http.StatusBadRequest, text)
}
func CreateForbiddenResponse(c *gin.Context, text string) {
	createResponse(c, http.StatusForbidden, text)
}

func CreateUnauthorizedResponse(c *gin.Context, text string) {
	createResponse(c, http.StatusUnauthorized, text)
}

func createResponse(c *gin.Context, statusCode int, error string) {
	c.AbortWithStatusJSON(statusCode, ErrorResponse{error})
}
