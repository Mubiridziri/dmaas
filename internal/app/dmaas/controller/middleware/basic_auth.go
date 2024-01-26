package middleware

import (
	"dmaas/internal/app/dmaas/controller/response"
	"dmaas/internal/app/dmaas/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const UserKey = "AUTH"

func AuthRequired(userRepository repository.UserRepositoryInterface) gin.HandlerFunc {

	return func(c *gin.Context) {
		session := sessions.Default(c)
		userKey := session.Get(UserKey)

		if userKey == nil {
			response.CreateUnauthorizedResponse(c, "unauthorized")
			return
		}

		user, err := userRepository.GetUserByUsername(userKey.(string))

		if err != nil {
			response.CreateUnauthorizedResponse(c, "invalid cookie")
			return
		}

		c.Set("user", user)
		c.Next()
	}

}
