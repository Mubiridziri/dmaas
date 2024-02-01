package v1

import (
	"dmaas/internal/context"
	"dmaas/internal/controller/response"
	"dmaas/internal/dto"
	"dmaas/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	Context *context.ApplicationContext
}

type PaginatedUsers struct {
	Total   int64         `json:"total"`
	Entries []entity.User `json:"entries"`
}

// listUsersAction GoDoc
//
//	@Summary	List User
//	@Schemes
//	@Description	Paginated User List
//	@Param			page	query	int	false	"Page"
//	@Param			limit	query	int	false	"Limit"
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	PaginatedUsers
//	@Router			/api/v1/users [GET]
func (controller *UserController) listUsersAction(c *gin.Context) {
	pagination, err := dto.QueryFromContext(c)

	if err != nil {
		response.CreateBadRequestResponse(c, err.Error())
	}

	entries, err := controller.Context.UserUseCase.ListUsers(pagination)
	count := controller.Context.UserUseCase.GetCount()

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, PaginatedUsers{
		Total:   count,
		Entries: entries,
	})
}

// createUserAction GoDoc
//
//	@Summary	Create User
//	@Schemes
//	@Description	Create entity
//	@Param			request	body	UserRequest	true	"User Data"
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.User
//	@Router			/api/v1/users [POST]
func (controller *UserController) createUserAction(c *gin.Context) {
	var request dto.UserRequest

	if err := c.BindJSON(&request); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	user := request.ToUser()
	err := controller.Context.UserUseCase.CreateUser(&user)

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, user)
}

// editUserAction GoDoc
//
//	@Summary	Update User
//	@Schemes
//	@Description	Update entity
//	@Tags			Users
//	@Param			id		path	int				true	"User ID"
//	@Param			request	body	UserRequest	true	"User Data"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.User
//	@Router			/api/v1/users/:id [PUT]
func (controller *UserController) editUserAction(c *gin.Context) {
	id, err := dto.IdFromContext(c)
	if err != nil {
		response.CreateNotFoundResponse(c, "invalid ID param")
		return
	}

	user, err := controller.Context.UserUseCase.GetUserById(id)
	if err != nil {
		response.CreateNotFoundResponse(c, "invalid ID param")
		return
	}

	var request dto.UserUpdateRequest
	if err := c.BindJSON(&request); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	request.ToUser(&user)
	err = controller.Context.UserUseCase.UpdateUser(&user)

	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// detailUserAction GoDoc
//
//	@Summary	Detail User
//	@Schemes
//	@Description	Get By ID
//	@Tags			Users
//	@Param			id	path	int	true	"User ID"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.User
//	@Router			/api/v1/users/:id [GET]
func (controller *UserController) detailUserAction(c *gin.Context) {
	id, err := dto.IdFromContext(c)
	if err != nil {
		response.CreateNotFoundResponse(c, "invalid ID param")
		return
	}

	user, err := controller.Context.UserUseCase.GetUserById(id)
	if err != nil {
		response.CreateNotFoundResponse(c, "not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

// removeUserAction GoDoc
//
//	@Summary	Remove User
//	@Schemes
//	@Description	Remove By ID
//	@Tags			Users
//	@Param			id	path	int	true	"User ID"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.User
//	@Router			/api/v1/users/:id [DELETE]
func (controller *UserController) removeUserAction(c *gin.Context) {
	id, err := dto.IdFromContext(c)
	if err != nil {
		response.CreateNotFoundResponse(c, "invalid ID param")
		return
	}

	user, err := controller.Context.UserUseCase.GetUserById(id)
	err = controller.Context.UserUseCase.RemoveUser(&user)
	if err != nil {
		response.CreateInternalServerResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (controller *UserController) AddUserRoute(r *gin.RouterGroup) {
	userGroup := r.Group("/users")

	userGroup.GET("", controller.listUsersAction)
	userGroup.POST("", controller.createUserAction)
	userGroup.PUT("/:id", controller.editUserAction)
	userGroup.GET("/:id", controller.detailUserAction)
	userGroup.DELETE("/:id", controller.removeUserAction)
}
