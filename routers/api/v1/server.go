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
package v1

import (
	"fmt"
	"net/http"

	"github.com/dmzlingyin/clipshare/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var transfer = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func server(c *gin.Context) {
	ws, err := transfer.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
	defer ws.Close()

	// 循环接收数据
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
		fmt.Println(mt, message)

		// 根据数据的发送方, 查找该发送方的所有在线设备, 并"多播"收到的数据
	}
}
