package middlewares

import (
	"github.com/N-SSL/container-target/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JudgeIsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		Role, _ := c.Get("role")
		if Role.(string) != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.Response{
				Success: false,
				Code:    http.StatusUnauthorized,
				Message: "You are not admin",
				Data:    nil,
			})
			return
		}
		c.Next()
	}
}
