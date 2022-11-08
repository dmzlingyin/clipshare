/*
* Get Token if the username is exist
* Author: Wenchao Lv
* Date: 2021-07-12
* Last modified: Wenchao Lv
* Last modified date: 2021-08-12, 2021-12-02
* Last modified content: migrate to RPA project
 */
package api

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/dmzlingyin/clipshare/pkg/app"
	"github.com/dmzlingyin/clipshare/pkg/e"
	"github.com/dmzlingyin/clipshare/pkg/model"
	"github.com/dmzlingyin/clipshare/pkg/utils"
	"github.com/gin-gonic/gin"
)

// use username and device as part of JWT PAYLOAD
type auth struct {
	uinfo
	Device string `valid:"Required; MaxSize(260)"`
}

func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	username := c.PostForm("username")
	password := c.PostForm("password")
	device := c.PostForm("device")
	a := auth{uinfo{username, password}, device}

	valid := validation.Validation{} // 实例化验证对象
	ok, _ := valid.Valid(&a)         // 验证参数是否符合约定
	if !ok {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	err := model.Check(username, password)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_USER_PASSWORD, nil)
		return
	}

	token, err := utils.GenerateToken(username, device)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}

func Auth(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
