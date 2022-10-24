//go:build client
// +build client

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	C "github.com/dmzlingyin/clipshare/pkg/constant"
	"github.com/dmzlingyin/clipshare/pkg/log"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
)

type Meta struct {
	UserName string
	Device   string
	Data     []byte
}

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "client can detect clipboard change and send the content to server",
	Run: func(cmd *cobra.Command, args []string) {
		client(C.ClientConf.UserName, C.ClientConf.Device)
	},
}

var ctx context.Context
var cancel context.CancelFunc

func init() {
	rootCmd.AddCommand(clientCmd)
	err := clipboard.Init()
	if err != nil {
		log.ErrorLogger.Fatalf("clipboard init failed")
	}

	// 开启剪贴板监控
	ctx, cancel = context.WithCancel(context.Background())
	go Watch(ctx, C.ClientConf.UserName, C.ClientConf.Device)
}

func client(username, device string) {
	fmt.Println(username, device, "connecting to server...")
	u := url.URL{Scheme: "ws", Host: C.ClientConf.Host, Path: "/socket"}

	// 建立websocket连接, 通过header区分客户端
	header := http.Header{}
	header.Add("UserName", username)
	header.Add("Device", device)
	c, _, err := (&websocket.Dialer{}).Dial(u.String(), header)
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.ErrorLogger.Println(err)
			return
		}
		if len(message) > 0 {
			// 关闭剪贴板监控
			cancel()
			fmt.Println(username, device, "recevied: ", string(message))
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
	u := url.URL{Scheme: "http", Host: C.ClientConf.Host, Path: "/transfer"}

	_, err := http.Post(u.String(), "application/json", bytes.NewReader(marshalForm))
	if err != nil {
		log.ErrorLogger.Println(err)
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
