package v1

import (
	"dmaas/internal/context"
	"dmaas/internal/controller/middleware"
	"dmaas/internal/controller/response"
	"dmaas/internal/entity"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SecurityController struct {
	Context *context.ApplicationContext
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginAction GoDoc
//
//	@Summary	Login
//	@Schemes
//	@Description	Authorization with help username and password
//	@Param			request	body	LoginRequest	true	"Username and password"
//	@Tags			Security
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	LoginRequest
//	@Router			/api/v1/login [POST]
func (controller *SecurityController) LoginAction(c *gin.Context) {
	session := sessions.Default(c)
	var request LoginRequest

	if err := c.BindJSON(&request); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	user, err := controller.Context.UserUseCase.GetUserByUsername(request.Username)

	if err != nil {
		response.CreateUnauthorizedResponse(c, "invalid credentials")
		return
	}

	if !user.IsPasswordCorrect(request.Password) {
		response.CreateUnauthorizedResponse(c, "invalid credentials")
		return
	}

	session.Set(middleware.UserKey, request.Username)
	_ = session.Save()

	c.JSON(http.StatusOK, request)

}

// ProfileAction GoDoc
//
//	@Summary	Profile
//	@Schemes
//	@Description	You can check auth or get profile data
//	@Tags			Security
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.User
//	@Router			/api/v1/login [GET]
func (controller *SecurityController) ProfileAction(c *gin.Context) {
	user := c.MustGet("user").(entity.User)
	c.JSON(http.StatusOK, user)
}

// LogoutAction GoDoc
//
//	@Summary	Logout
//	@Schemes
//	@Description	Logout from account
//	@Tags			Security
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/api/v1/logout [POST]
func (controller *SecurityController) LogoutAction(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(middleware.UserKey)
	err := session.Save()
	if err != nil {
		response.CreateInternalServerResponse(c, "error save session!")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
