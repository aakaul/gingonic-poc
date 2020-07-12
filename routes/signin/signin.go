package signin

import (
	"authPoc/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//SignIn ...
type SignIn interface {
	save(models.User) models.User
	sigin() models.User
}

func sigin(c *gin.Context) {

	session := sessions.Default(c)

	sessionID := session.Get("uid")

	if sessionID != nil {
		c.Redirect(http.StatusFound, "/")
	} else {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "Logged In",
		})
	}

}

func signIn(c *gin.Context) {
	session := sessions.Default(c)

	var user models.User
	c.BindJSON(&user)

	if user.Password == "" {
		c.JSON(401, gin.H{
			"msg": "need all params",
			"err": "unauth",
		})
		return
	}

	uname, err := models.GetKeys("users:*:" + user.Uname)
	if err != nil || len(uname) == 0 {
		c.JSON(401, gin.H{
			"msg": "something went wrong",
			"err": "unauth",
		})
		panic(err)
	}

	password, _ := models.Rdb.HGet(uname[0], "password").Result()
	if password == user.Password {
		session.Set("uid", user.UID)
		session.Save()
		c.JSON(200, gin.H{
			"data": user,
		})
		return
	}
	c.JSON(401, gin.H{
		"msg": "something went wrong",
		"err": "unauth",
	})
	return

}

func register(c *gin.Context) {
	session := sessions.Default(c)
	sessionID := session.Get("uid")

	if sessionID != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	var user map[string]interface{}
	c.BindJSON(&user)

	if user["password"] == "" || user["uname"] == "" {
		c.JSON(401, gin.H{
			"msg": "need all params",
			"err": "unauth",
		})
		panic("not found")
	}
	uid := models.WriteInDBAndRdb(user, "users")
	session.Set("uid", uid)
	session.Save()
	c.Redirect(http.StatusFound, "/")

}
