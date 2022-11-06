package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Response struct {
	Success bool        `json:"success" example:"true"`
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"example message"`
	Data    interface{} `json:"data"`
}

func BadRequest(c *gin.Context) {
	log.Println("Bad request")
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Code:    http.StatusBadRequest,
		Message: "Bad request.",
		Data:    nil,
	})
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

