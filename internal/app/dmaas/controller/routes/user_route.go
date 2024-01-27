package routes

import (
	"dmaas/internal/app/dmaas/controller/response"
	"dmaas/internal/app/dmaas/entity"
	"dmaas/internal/app/dmaas/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserController struct {
	Repository repository.UserRepositoryInterface
}

type UserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
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
//	@Param			limit	query	int	false	"Page"
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	PaginatedUsers
//	@Router			/api/v1/users [GET]
func (controller *UserController) listUsersAction(c *gin.Context) {
	//TODO may be bind to model (struct) ?
	pageQuery := c.DefaultQuery("page", "1")
	limitQuery := c.DefaultQuery("limit", "10")

	page, pageOk := strconv.Atoi(pageQuery)
	limit, limitOk := strconv.Atoi(limitQuery)

	if pageOk != nil || limitOk != nil {
		response.CreateBadRequestResponse(c, "bad query parameters")
		return
	}

	entries, err := controller.Repository.ListUsers(page, limit)
	count := controller.Repository.GetCount()

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
	var request UserRequest

	if err := c.BindJSON(&request); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	//TODO MOVE TO .... ????
	//boilerplate
	user := entity.User{
		Name:     request.Name,
		Username: request.Username,
		Password: request.Password,
	}
	err := controller.Repository.CreateUser(&user)

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
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		response.CreateNotFoundResponse(c, "invalid ID param")
		return
	}

	user, err := controller.Repository.GetUserById(id)

	if err != nil {
		response.CreateNotFoundResponse(c, "invalid ID param")
		return
	}

	var request UserRequest

	if err := c.BindJSON(&request); err != nil {
		response.CreateBadRequestResponse(c, err.Error())
		return
	}

	//TODO Refact?? need object to populate method
	//boilerplate
	user.Name = request.Name
	user.Username = request.Username
	user.Password = request.Password

	err = controller.Repository.UpdateUser(&user)

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
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		response.CreateNotFoundResponse(c, "invalid ID param")
		return
	}

	user, err := controller.Repository.GetUserById(id)

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
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		response.CreateNotFoundResponse(c, "invalid ID param")
		return
	}

	user, err := controller.Repository.GetUserById(id)
	err = controller.Repository.RemoveUser(&user)

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
