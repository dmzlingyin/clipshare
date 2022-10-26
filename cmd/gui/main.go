//go:build darwin || linux || windows
// +build darwin linux windows

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"golang.design/x/clipboard"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

var ctx context.Context
var cancel context.CancelFunc

type Meta struct {
	UserName string
	Device   string
	Data     []byte
}

const (
	username = "lingyin"
	device   = "android"
	host     = "172.17.130.166:8080"
)

func init() {
	err := clipboard.Init()
	if err != nil {
		return
	}
}

func main() {
	go connect()
	Watch(ctx, username, device)
	app.Main(func(a app.App) {
		var glctx gl.Context
		for {
			select {
			case e := <-a.Events():
				switch e := a.Filter(e).(type) {
				case lifecycle.Event:
					switch e.Crosses(lifecycle.StageAlive) {
					case lifecycle.CrossOff:
						send("android", "lifecycle stageAlive", []byte("crossOff"))
					case lifecycle.CrossOn:
						send("android", "lifecycle stageAlive", []byte("crossOn"))
					}
					switch e.Crosses(lifecycle.StageDead) {
					case lifecycle.CrossOff:
						send("android", "lifecycle stageDead", []byte("crossOff"))
					case lifecycle.CrossOn:
						send("android", "lifecycle stageDead", []byte("crossOn"))
					}
					switch e.Crosses(lifecycle.StageVisible) {
					case lifecycle.CrossOff:
						send("android", "lifecycle stageVisible", []byte("crossOff"))
					case lifecycle.CrossOn:
						send("android", "lifecycle stageVisible", []byte("crossOn"))
					}
					switch e.Crosses(lifecycle.StageFocused) {
					case lifecycle.CrossOff:
						send("android", "lifecycle stageFocused", []byte("crossOff"))
					case lifecycle.CrossOn:
						send("android", "lifecycle stageFocused", []byte("crossOn"))
					}
					glctx, _ = e.DrawContext.(gl.Context)
				case touch.Event:
					send("android", "touch enent", []byte("touched"))
				case paint.Event:
					if glctx == nil {
						continue
					}
					onDraw(glctx)
					a.Publish()
				}
			}
		}
	})
}

var (
	determined = make(chan struct{})
	ok         = false
)

func connect() {
	u := url.URL{Scheme: "ws", Host: host, Path: "/socket"}

	// 建立websocket连接, 通过header区分客户端
	header := http.Header{}
	header.Add("UserName", "lingyin")
	header.Add("Device", "android")
	c, _, err := (&websocket.Dialer{}).Dial(u.String(), header)
	defer c.Close()
	if err != nil {
		return
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return
		}
		if len(message) > 0 {
			clipboard.Write(clipboard.FmtText, message)
		}
	}
}

func onDraw(glctx gl.Context) {
	glctx.ClearColor(0, 1, 0, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
}

func send(username, device string, data []byte) {
	form := Meta{
		UserName: username,
		Device:   device,
		Data:     data,
	}
	marshalForm, _ := json.Marshal(form)
	u := url.URL{Scheme: "http", Host: host, Path: "/transfer"}

	_, err := http.Post(u.String(), "application/json", bytes.NewReader(marshalForm))
	if err != nil {
		return
	}
}

func Watch(ctx context.Context, username, device string) {
	ti := time.NewTicker(time.Second)
	last := clipboard.Read(clipboard.FmtText)
	go func() {
		for {
			select {
			case <-ti.C:
				b := clipboard.Read(clipboard.FmtText)
				if b == nil {
					continue
				}
				if !bytes.Equal(last, b) {
					send(username, device, b)
					last = b
				}
			}
		}
	}()
}
