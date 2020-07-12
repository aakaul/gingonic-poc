package signin

import (
	"github.com/gin-gonic/gin"
)

//signin.SignInRoutes
func SignInRoutes(route *gin.RouterGroup) {
	signin := route.Group("/signin")
	{
		signin.GET("/", sigin)
		// signin.GET("/:id", )
		signin.POST("/", signIn)
	}

	route.GET("/signOut", signOut)
	route.POST("/register", register)
}
