package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID := session.Get("uid")
		fmt.Print(sessionID)
		if sessionID == nil {
			c.Redirect(http.StatusFound, "/api/signin")
			c.Abort()
		}
	}
}
