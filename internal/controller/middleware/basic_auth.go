package middleware

import (
	"dmaas/internal/context"
	"dmaas/internal/controller/response"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const UserKey = "AUTH"

func AuthRequired(context *context.ApplicationContext) gin.HandlerFunc {

	return func(c *gin.Context) {
		session := sessions.Default(c)
		userKey := session.Get(UserKey)

		if userKey == nil {
			response.CreateUnauthorizedResponse(c, "unauthorized")
			return
		}

		user, err := context.UserUseCase.GetUserByUsername(userKey.(string))

		if err != nil {
			response.CreateUnauthorizedResponse(c, "invalid cookie")
			return
		}

		c.Set("user", user)
		c.Next()
	}

}
