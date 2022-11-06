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

package v1

import (
	"net/http"
	"time"

	C "github.com/dmzlingyin/clipshare/pkg/constant"
	"github.com/dmzlingyin/clipshare/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type conn struct {
	device string
	ws     *websocket.Conn
}

type UD struct {
	UserName string `form:"username"`
	Device   string `form:"device"`
}

type Meta struct {
	UD
	Data []byte
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

// Socket 与客户端建立连接
func Socket(c *gin.Context) {
	var ud UD
	c.Bind(&ud)

	// 超过服务器的最大允许用户数量
	if len(conns) > C.ServerConf.MaxUsers {
		log.WarningLogger.Printf("user %s's connecting was refused\n", ud.UserName)
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "too many users connected",
		})
		return
	}
	// 超过服务器的最大允许用户设备数量
	if len(conns[ud.UserName]) > C.ServerConf.MaxDevices {
		log.WarningLogger.Printf("user %s's device %s connecting was refused\n", ud.UserName, ud.Device)
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "too many devices connected",
		})
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.ErrorLogger.Println(err)
	}
	defer ws.Close()

	// 将ws加入对应用户的连接队列
	conns[ud.UserName] = append(conns[ud.UserName], conn{device: ud.Device, ws: ws})
	log.InfoLogger.Println(ud.UserName, ud.Device, "online")

	// 心跳检测(1s)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		err := ws.WriteMessage(websocket.TextMessage, []byte{})
		if err != nil {
			log.ErrorLogger.Println(ud.UserName, ud.Device, "offline")
			// 从连接队列移除
			q := conns[ud.UserName]
			for i := 0; i < len(q); i++ {
				if q[i].device == ud.Device {
					q[i].ws.Close()
					conns[ud.UserName] = append(conns[ud.UserName][:i], conns[ud.UserName][i+1:]...)
				}
			}
			return
		}
	}
}

// Transfer接口获取传入数据, 并进行广播
func Transfer(c *gin.Context) {
	// 读取发送用户、device、数据信息
	var userInfo Meta
	c.Bind(&userInfo)

	log.InfoLogger.Println(userInfo.UserName, userInfo.Device, "sended: ", string(userInfo.Data))
	// 向发送方其他在线设备进行广播
	for _, conn := range conns[userInfo.UserName] {
		if conn.device != userInfo.Device {
			err := conn.ws.WriteMessage(websocket.TextMessage, userInfo.Data)
			if err != nil {
				log.ErrorLogger.Println("data send to", userInfo.UserName, userInfo.Data, "error")
				continue
			}
		}
	}
}
