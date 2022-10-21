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
	"time"

	"github.com/dmzlingyin/clipshare/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type conn struct {
	device string
	ws     *websocket.Conn
}

var (
	// 维护每个用户的连接队列
	conns    = map[string][]conn{}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

func Socket(c *gin.Context) {
	fmt.Println(c.Request.Header["Username"])
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.ErrorLogger.Println(err)
	}
	defer ws.Close()

	// 将ws加入对应用户的连接队列
	conns["username"] = append(conns["username"], conn{device: "test", ws: ws})

	// 心跳检测(1s)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		fmt.Println("1s done")
		err := ws.WriteMessage(websocket.TextMessage, []byte{'\x00'})
		if err != nil {
			log.ErrorLogger.Println("device offline")
			return
		}
	}
}

// func Transfer(c *gin.Context) {
// 	// 读取用户、device_id、数据信息
// 	username := "test"
// 	device := "test"
// 	data := "test"

// 	for conn := range conns["username"] {
// 		if conn.device != device {
// 			conn.ws.WriteMessage()
// 		}
// 	}
// }
