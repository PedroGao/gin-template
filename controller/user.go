package controller

import (
	"fmt"
	"github.com/PedroGao/jerry/form"
	"github.com/PedroGao/jerry/libs/erro"
	"github.com/PedroGao/jerry/libs/token"
	"github.com/PedroGao/jerry/model"
	"github.com/PedroGao/jerry/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	var (
		register form.Register
		err      error
	)
	if err = c.ShouldBindJSON(&register); err != nil {
		c.Error(err)
		return
	}
	_, err = model.CreateUser(register.Nickname, register.Password)
	if err != nil {
		c.Error(erro.ParamsErr.SetMsg(err.Error()))
	}
	c.JSON(http.StatusOK, erro.OK)
}

func Login(c *gin.Context) {
	var (
		login    form.Login
		tokenStr string
		err      error
	)
	// ATTENTION: 不要用MustBind，它会自动将 Content-Type 设置成 text/plain
	if err = c.ShouldBindJSON(&login); err != nil {
		c.Error(err)
		return
	}

	// COMMAND: 在controller里面使用erro，在其他的地方使用 errors.New等原生的error
	if err = login.ValidateNameAndPassword(); err != nil {
		c.Error(erro.ParamsErr.SetMsg(err.Error()))
		return
	}

	tokenStr, err = token.GenerateAccessToken(login.Nickname)

	if err != nil {
		c.Error(erro.ParamsErr.SetMsg(err.Error()))
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenStr,
	})
}

func GetUsers(c *gin.Context) {
	infos, e := service.ListUser()
	value, exists := c.Get("user")
	if exists {
		fmt.Println(value)
	}
	if e != nil {
		c.Error(erro.UserNotFound)
		return
	}
	c.JSON(http.StatusOK, infos)
}
