package middleware

import (
	"dmaas/internal/app/dmaas/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const UserKey = "AUTH"

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	userKey := session.Get(UserKey)

	if userKey == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user, err := repository.FindUserByUsername(userKey.(string))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid cookie",
		})
		return
	}

	c.Set("user", user)
	c.Next()
}
