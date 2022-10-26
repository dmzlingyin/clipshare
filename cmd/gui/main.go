//go:build darwin || linux || windows
// +build darwin linux windows

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

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
	ctx    gl.Context
	header http.Header
}

func NewClipApp(a app.App) *ClipApp {
	header := http.Header{}
	header.Add("username", username)
	header.Add("device", device)
	return &ClipApp{app: a, header: header}
}

func (c *ClipApp) Watch(format clipboard.Format) {
	t := time.NewTicker(time.Second)
	last := clipboard.Read(format)
	for range t.C {
		data := clipboard.Read(format)
		if len(data) == 0 {
			continue
		}
		if !bytes.Equal(last, data) {
			c.Send(data)
			last = data
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
			clipboard.Write(clipboard.FmtText, message)
		}
	}
}

func (c *ClipApp) draw() {
	c.ctx.ClearColor(0, 1, 0, 1)
	c.ctx.Clear(gl.COLOR_BUFFER_BIT)
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
		clip.Connect()
		clip.Watch(clipboard.FmtText)

		for e := range clip.app.Events() {
			switch e := clip.app.Filter(e).(type) {
			case paint.Event:
				clip.draw()
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					clip.Watch(clipboard.FmtImage)
				}
			}
		}
	})
}
