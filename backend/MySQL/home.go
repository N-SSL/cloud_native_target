package MySQL

import (
	"github.com/N-SSL/container-target/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPersonalInfo(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.Response{
			Success: false,
			Code:    http.StatusUnauthorized,
			Message: "get user info error",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, controllers.Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    user,
	})
}
