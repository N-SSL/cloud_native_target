package MySQL

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/N-SSL/container-target/controllers"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func JustRegister(c *gin.Context)  {
	var registerForm Login
	if err := c.Bind(&registerForm); err != nil {
		controllers.BadRequest(c)
		return
	}
	bytepwd,_ := bcrypt.GenerateFromPassword([]byte(registerForm.Password),0)
	var newUser = Users{
		Username: registerForm.Username,
		Password: bytepwd,
		Role: "normal",
	}

	if err := SqlDB.Create(&newUser).Error; err!=nil {
		log.Println("cannot register.")
		c.AbortWithStatusJSON(http.StatusBadRequest, controllers.Response{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "cannot register.",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK,controllers.Response{
		Success: true,
		Code: http.StatusOK,
		Message: "register success",
		Data: nil,
	})
}

func getMD5Hash(TestString string) string {
	Md5Inst := md5.New()
	Md5Inst.Write([]byte(TestString))
	Result := hex.EncodeToString(Md5Inst.Sum(nil))
	return Result
}