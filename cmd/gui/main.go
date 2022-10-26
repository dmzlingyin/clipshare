//go:build darwin || linux || windows
// +build darwin linux windows

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
	ctx    context.Context
	cancel context.CancelFunc
}

func NewClipApp(a app.App) *ClipApp {
	header := http.Header{}
	header.Add("username", username)
	header.Add("device", device)
	ctx, cancel := context.WithCancel(context.Background())
	return &ClipApp{app: a, header: header, ctx: ctx, cancel: cancel}
}

func (c *ClipApp) Watch() {
	ch := clipboard.Watch(c.ctx, clipboard.FmtText)
	for data := range ch {
		if len(data) > 0 {
			c.Send(data)
		}
	}
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

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			return
		}
		if len(message) > 0 {
			c.cancel()
			clipboard.Write(clipboard.FmtText, message)
			c.ctx, c.cancel = context.WithCancel(context.Background())
			go c.Watch()
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
	c.gctx.ClearColor(0, 1, 0, 1)
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
		go clip.Watch()

		for e := range clip.app.Events() {
			switch e := clip.app.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOff:
					clip.OnStop()
				case lifecycle.CrossOn:
					clip.OnStart(e)
				}
			case paint.Event:
				clip.draw()
			}
		}
	})
}
