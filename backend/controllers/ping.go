package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping(c *gin.Context)  {
		text := c.Param("text")
		c.String(http.StatusOK,text)
}
