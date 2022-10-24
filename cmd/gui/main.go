//go:build android
// +build android

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"golang.design/x/clipboard"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/gl"
)

const (
	username = "lingyin"
	device   = "android"
)

type Meta struct {
	UserName string
	Device   string
	Data     []byte
}

var data []byte
var count = 1
var ch chan bool
var ctx context.Context
var cancel context.CancelFunc

func init() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	ctx, cancel = context.WithCancel(context.TODO())
	go Watch(ctx, username, device)
}

func main() {
	go connect()

	app.Main(func(a app.App) {
		var glctx gl.Context
		for {
			select {
			case <-ch:
				a.Send(paint.Event{})

			case e := <-a.Events():
				switch e := a.Filter(e).(type) {
				case lifecycle.Event:
					glctx, _ = e.DrawContext.(gl.Context)
				case paint.Event:
					if glctx == nil {
						continue
					}
					draw(glctx)
					a.Publish()
				}
			}
		}
	})
}

func draw(glctx gl.Context) {
	if count%2 == 0 {
		glctx.ClearColor(0, 1, 0, 1)
	} else {
		glctx.ClearColor(0, 1, 1, 1)
	}
	count++
	glctx.Clear(gl.COLOR_BUFFER_BIT)
}

func connect() {
	u := url.URL{Scheme: "ws", Host: "172.17.130.166:8080", Path: "/socket"}

	// 建立websocket连接, 通过header区分客户端
	header := http.Header{}
	header.Add("UserName", username)
	header.Add("Device", device)
	c, _, err := (&websocket.Dialer{}).Dial(u.String(), header)
	if err != nil {
		return
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return
		}
		if len(message) > 0 {
			// 关闭剪贴板监控
			cancel()
			ch <- true
			clipboard.Write(clipboard.FmtText, message)
			// 重新开启剪贴板监控
			ctx, cancel = context.WithCancel(context.Background())
			go Watch(ctx, username, device)
		}
	}
}

func send(username, device string, data []byte) {
	form := Meta{
		UserName: username,
		Device:   device,
		Data:     data,
	}
	marshalForm, _ := json.Marshal(form)
	u := url.URL{Scheme: "http", Host: "172.17.130.166:8080", Path: "/transfer"}

	_, err := http.Post(u.String(), "application/json", bytes.NewReader(marshalForm))
	if err != nil {
		return
	}
}

func Watch(ctx context.Context, username, device string) {
	ch := clipboard.Watch(ctx, clipboard.FmtText)
	for data := range ch {
		if len(data) != 0 {
			send(username, device, data)
		}
	}
}
