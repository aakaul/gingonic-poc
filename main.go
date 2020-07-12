package main

import (
	"authPoc/models"
	"authPoc/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	app := gin.Default()

	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	app.Use(sessions.Sessions("userSessions", store))

	app.Static("/assets", "./public")
	models.ConnectDataBase()
	routes.Router(app)

	app.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
