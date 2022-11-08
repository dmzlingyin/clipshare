//go:build client
// +build client

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/dmzlingyin/clipshare/hub"
	C "github.com/dmzlingyin/clipshare/pkg/constant"
	"github.com/dmzlingyin/clipshare/pkg/e"
	"github.com/dmzlingyin/clipshare/pkg/log"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

type rv struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

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

	// 登录验证
	login()
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

func login() {
	u := url.URL{Scheme: "http", Host: C.ClientConf.Host, Path: "/login"}
	v := url.Values{}

	if C.ClientConf.Token == "" {
		v.Set("username", C.ClientConf.UserName)
		v.Set("password", C.ClientConf.PassWord)
		v.Set("device", C.ClientConf.Device)
		resp, err := http.PostForm(u.String(), v)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.ErrorLogger.Fatal(err)
		}

		rvalue := rv{}
		err = json.Unmarshal(body, &rvalue)
		if err != nil {
			log.ErrorLogger.Fatal(err)
		}
		if rvalue.Code == e.ERROR_USER_PASSWORD {
			// 自动注册
			u.Path = "/register"
			resp, err = http.PostForm(u.String(), v)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			body, err = io.ReadAll(resp.Body)
			err = json.Unmarshal(body, &rvalue)
			if err != nil {
				panic(err)
			}
		}

		// 获取的token写入文件
		fmt.Println("请重新执行程序")
		os.Exit(0)
	} else {
		u.Path = "/api/v1/auth"
		client := &http.Client{}
		req, _ := http.NewRequest("GET", u.String(), nil)
		req.Header.Add("Token", C.ClientConf.Token)
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(body, &rvalue)

		if rvalue.Code != http.StatusOK {
			log.ErrorLogger.Println("token invalid")
			// 如果token不合法, 重新获取token
			C.ClientConf.Token = ""
			login()
		}
	}
}
