package middleware

import (
	"dmaas/internal/app/dmaas/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const UserKey = "AUTH"

func AuthRequired(userRepository repository.UserRepositoryInterface) gin.HandlerFunc {

	return func(c *gin.Context) {
		session := sessions.Default(c)
		userKey := session.Get(UserKey)

		if userKey == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		user, err := userRepository.GetUserByUsername(userKey.(string))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid cookie",
			})
			return
		}

		c.Set("user", user)
		c.Next()
	}

}
