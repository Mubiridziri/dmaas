package routes

import (
	"dmaas/internal/app/dmaas/controller/middleware"
	"dmaas/internal/app/dmaas/controller/response"
	"dmaas/internal/app/dmaas/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginAction(c *gin.Context) {
	session := sessions.Default(c)
	var request LoginRequest

	if err := c.BindJSON(&request); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	user, err := repository.FindUserByUsername(request.Username)

	if err != nil {
		response.CreateUnauthorizedResponse(c, "Invalid credentials")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		response.CreateUnauthorizedResponse(c, "Invalid credentials")
		return
	}

	session.Set(middleware.UserKey, request.Username)
	_ = session.Save()

	c.JSON(http.StatusAccepted, request)

}
func LogoutAction(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(middleware.UserKey)
	err := session.Save()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Error save session!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
