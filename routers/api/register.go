package api

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/dmzlingyin/clipshare/pkg/app"
	"github.com/dmzlingyin/clipshare/pkg/e"
	"github.com/dmzlingyin/clipshare/pkg/model"
	"github.com/gin-gonic/gin"
)

type uinfo struct {
	UserName string `valid:"Required; MaxSize(260)"`
	PassWord string `valid:"Required; MaxSize(260)"`
}

func Register(c *gin.Context) {
	appG := app.Gin{C: c}
	username := c.PostForm("username")
	password := c.PostForm("password")

	a := uinfo{username, password}
	valid := validation.Validation{} // 实例化验证对象
	ok, _ := valid.Valid(&a)         // 验证参数是否符合约定
	if !ok {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	err := model.Register(username, password)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_USER_CREATE, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
