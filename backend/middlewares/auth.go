package middlewares

import (
	"github.com/N-SSL/container-target/MySQL"
	"github.com/N-SSL/container-target/controllers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		id := session.Get("id")
		if id == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.Response{
				Success: false,
				Code:    http.StatusUnauthorized,
				Message: "Please login.",
				Data:    nil,
			})
			return
		}
		var user MySQL.Users
		// MySQL.SqlDB.Where(id).Preload("Roles").First(&user)
		// SqlDB.Where(&MySQL.Users{ID: id}).First(&user)
		MySQL.SqlDB.Where(id).First(&user)
		c.Set("user", user)
		c.Set("id", strconv.FormatInt(id.(int64),10))
		c.Set("role", user.Role)
		c.Next()
	}
}
