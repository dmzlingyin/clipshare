//go:build server
// +build server

package jwt

import (
	"net/http"
	"time"

	"github.com/dmzlingyin/clipshare/pkg/app"
	"github.com/dmzlingyin/clipshare/pkg/e"
	"github.com/dmzlingyin/clipshare/pkg/utils"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		code := e.SUCCESS

		token, ok := c.Request.Header["Token"]
		if !ok {
			code = e.ERROR_AUTH
			appG.Response(http.StatusBadRequest, code, nil)
			c.Abort()
			return
		}

		if token[0] == "" {
			code = e.ERROR_AUTH
		} else {
			claims, err := utils.ParseToken(token[0])
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
			if err == nil {
				c.Set("username", claims.UserName)
				c.Set("device", claims.Device)
			}
		}

		if code != e.SUCCESS {
			appG.Response(http.StatusUnauthorized, code, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
