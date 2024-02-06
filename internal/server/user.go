package server

import (
	"dmaas/internal/usecase/users"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// handleLogin GoDoc
//
//	@Summary	Login
//	@Schemes
//	@Description	Authorization with help username and password
//	@Param			request	body	users.UserLogin	true	"Model"
//	@Tags			Security
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	users.UserView
//	@Router			/api/v1/login [POST]
func (s *Server) handleLogin(c *gin.Context) {
	var login users.UserLogin

	if err := c.ShouldBindJSON(&login); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := s.userController.LoginUser(login)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	session := sessions.Default(c)
	session.Set(UserKey, user.Username)
	err = session.Save()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "cannot save cookies",
		})
		return
	}

	c.JSON(http.StatusOK, user)

}

// handleLogout GoDoc
//
//	@Summary	Logout
//	@Schemes
//	@Description	Logout from account
//	@Tags			Security
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/api/v1/logout [POST]
func (s *Server) handleLogout(c *gin.Context) {
	user := c.MustGet("user")

	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	session := sessions.Default(c)
	session.Delete(UserKey)
	err := session.Save()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "cannot save cookies",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// handleProfile GoDoc
//
//	@Summary	Profile
//	@Schemes
//	@Description	You can check auth or get profile data
//	@Tags			Security
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	users.UserView
//	@Router			/api/v1/login [GET]
func (s *Server) handleProfile(c *gin.Context) {
	user := c.MustGet("user")

	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) handleListUser(c *gin.Context) {
	var query ListQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad query params: " + err.Error(),
		})
		return
	}

	rows, err := s.userController.ListUsers(query.Page, query.Limit)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Server error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, rows)

}

func (s *Server) handleCreateUser(c *gin.Context) {
	var createUser users.CreateOrUpdateUserView

	if err := c.ShouldBindJSON(&createUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Bind error: " + err.Error(),
		})
		return
	}

	user, err := s.userController.CreateUser(createUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Server error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)

}

func (s *Server) handleUpdateUser(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID param in path",
		})
		return
	}

	var createUser users.CreateOrUpdateUserView

	if err := c.ShouldBindJSON(&createUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Bind error: " + err.Error(),
		})
		return
	}

	user, err := s.userController.UpdateUser(id, createUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Server error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) handleDetailUser(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID param in path",
		})
		return
	}

	user, err := s.userController.GetUserById(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) handleDeleteUser(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID param in path",
		})
		return
	}

	user, err := s.userController.RemoveUser(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) AddUserRoutes(g *gin.RouterGroup) {
	grp := g.Group("/users")
	grp.GET("", s.handleListUser)
	grp.POST("", s.handleCreateUser)
	grp.PUT("/:id", s.handleUpdateUser)
	grp.GET("/:id", s.handleDetailUser)
	grp.DELETE("/:id", s.handleDeleteUser)
}
