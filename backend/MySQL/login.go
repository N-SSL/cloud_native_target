package MySQL

import (
	"github.com/N-SSL/container-target/controllers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func JustLogin(c *gin.Context)  {
	var loginForm Login
	if err := c.Bind(&loginForm); err != nil {
		controllers.BadRequest(c)
		return
	}
	var user Users
	if err := user.auth(loginForm); err != nil {
		c.JSON(http.StatusUnauthorized, controllers.Response{
			Success: false,
			Code:    http.StatusUnauthorized,
			Message: "Wrong email or password",
			Data:    nil,
		})
		return
	}
	session := sessions.Default(c)
	session.Set("id", user.ID) // set cookie and session
	if err := session.Save(); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, controllers.Response{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "server error",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, controllers.Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "Login successfully.",
		Data:    user,
	})
}

func (user *Users) auth(loginForm Login) error {
	if err := SqlDB.Where(&Users{Username: loginForm.Username}).First(&user).Error; err != nil {
		return err
	}
	log.Println(user.Password)
	log.Println(loginForm.Password)
	if err := bcrypt.CompareHashAndPassword(user.Password,[]byte(loginForm.Password)); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
