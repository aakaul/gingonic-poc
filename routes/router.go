package routes

import (
	"authPoc/middleware"
	"authPoc/routes/signin"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.LoadHTMLGlob("views/*")

	r.GET("/", middleware.AuthRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Logged In",
		})
	})

	user := r.Group("/api")
	{
		signin.SignInRoutes(user)
	}
}
