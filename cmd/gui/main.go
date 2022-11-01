//go:build darwin || linux || windows
// +build darwin linux windows

package main

import (
	"bytes"
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
	host     = "172.17.130.166:8080"
)

type Meta struct {
	UserName string
	Device   string
	Data     []byte
}

type ClipApp struct {
	app    app.App
	gctx   gl.Context
	header http.Header
}

var conncted = false

func NewClipApp(a app.App) *ClipApp {
	header := http.Header{}
	header.Add("username", username)
	header.Add("device", device)
	return &ClipApp{app: a, header: header}
}

func (c *ClipApp) Send(data []byte) {
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

func (c *ClipApp) Connect() {
	u := url.URL{Scheme: "ws", Host: host, Path: "/socket"}
	// 建立websocket连接, 通过header区分客户端
	ws, _, err := (&websocket.Dialer{}).Dial(u.String(), c.header)
	if err != nil {
		return
	}
	defer ws.Close()
	conncted = true

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			return
		}
		if len(message) > 0 {
			clipboard.Write(clipboard.FmtText, message)
		}
	}
}

func (c *ClipApp) OnStart(e lifecycle.Event) {
	c.gctx, _ = e.DrawContext.(gl.Context)
	c.app.Send(paint.Event{})
}

func (c *ClipApp) OnStop() {
	c.gctx = nil
}

func (c *ClipApp) draw() {
	if c.gctx == nil {
		return
	}
	defer c.app.Send(paint.Event{})
	defer c.app.Publish()
	if conncted {
		c.gctx.ClearColor(0, 1, 0, 1)
	} else {
		c.gctx.ClearColor(1, 0, 0, 1)
	}
	c.gctx.Clear(gl.COLOR_BUFFER_BIT)
}

func init() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
}

func main() {
	app.Main(func(a app.App) {
		clip := NewClipApp(a)
		clip.app.Send(paint.Event{})
		go clip.Connect()

		for e := range clip.app.Events() {
			switch e := clip.app.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOff:
					clip.OnStop()
				case lifecycle.CrossOn:
					clip.OnStart(e)
					data := clipboard.Read(clipboard.FmtText)
					clip.Send(data)
				}
			case paint.Event:
				clip.draw()
			}
		}
	})
}
