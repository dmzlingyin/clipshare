//go:build client
// +build client

package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dmzlingyin/clipshare/hub"
	C "github.com/dmzlingyin/clipshare/pkg/constant"
	"github.com/dmzlingyin/clipshare/pkg/log"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "client can detect clipboard change and send the content to server",
	Run: func(cmd *cobra.Command, args []string) {
		client()
	},
}

// 用于控制剪贴板监控开启或关闭
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
	go watch(ctx, C.ClientConf.UserName, C.ClientConf.Device)
	// 注册热键
	go mainthread.Init(fn)
}

func fn() {
	// Ctrl + Shift + S
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
	err := hk.Register()
	if err != nil {
		log.ErrorLogger.Println("热键注册失败")
		log.ErrorLogger.Println(err)
	}

	for {
		<-hk.Keydown()
		// 文件分享
		fmt.Println("hotkey down")
	}
}

func client() {
	u := url.URL{Scheme: "ws", Host: C.ClientConf.Host, Path: "/socket"}

	// 建立websocket连接, 通过header区分客户端
	header := http.Header{}
	header.Add("UserName", C.ClientConf.UserName)
	header.Add("Password", C.ClientConf.PassWord)
	header.Add("Device", C.ClientConf.Device)
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
			clipboard.Write(clipboard.FmtText, message)
			// 重新开启剪贴板监控
			ctx, cancel = context.WithCancel(context.Background())
			go watch(ctx, C.ClientConf.UserName, C.ClientConf.Device)
		}
	}
}

func watch(ctx context.Context, username, device string) {
	ch := clipboard.Watch(ctx, clipboard.FmtText)
	for data := range ch {
		if len(data) != 0 {
			hub.Send(username, device, data)
		}
	}
}
