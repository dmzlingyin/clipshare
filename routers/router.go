//go:build server
// +build server

/*
Copyright © 2022 Whenchao Lv dmzlingyin@163.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package routers

import (
	"github.com/dmzlingyin/clipshare/middleware/jwt"
	"github.com/dmzlingyin/clipshare/routers/api"
	v1 "github.com/dmzlingyin/clipshare/routers/api/v1"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	// 允许跨域请求
	r.Use(cors.Default())

	r.POST("/login", api.GetAuth)
	apiv1 := r.Group("/api/v1")

	apiv1.Use(jwt.JWT())
	{
		r.GET("/socket", v1.Socket)
		r.POST("/transfer", v1.Transfer)
	}
	return r
}
