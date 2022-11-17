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
	"time"

	"github.com/dmzlingyin/clipshare/hub"
	C "github.com/dmzlingyin/clipshare/pkg/constant"
	"github.com/dmzlingyin/clipshare/pkg/e"
	"github.com/dmzlingyin/clipshare/pkg/log"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
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
	go watch(ctx, C.Token)
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
	u := url.URL{Scheme: "ws", Host: C.ClientConf.Host, Path: "/api/v1/socket"}

	// 建立websocket连接, 通过header中的token区分用户和客户端
	header := http.Header{}
	header.Add("Token", C.Token)
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
			go watch(ctx, C.Token)
		}
	}
}

func watch(ctx context.Context, token string) {
	ch := clipboard.Watch(ctx, clipboard.FmtText)
	for data := range ch {
		if len(data) != 0 {
			err := hub.Send(token, data)
			if err != nil {
				log.ErrorLogger.Println(err)
			} else {
				// 播放成功提示音
				f, err := os.Open("./docs/send.wav")
				streamer, format, err := wav.Decode(f)
				if err != nil {
					log.ErrorLogger.Println(err)
				}
				defer streamer.Close()
				speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
				speaker.Play(streamer)
			}
		}
	}
}

func login() {
	u := url.URL{Scheme: "http", Host: C.ClientConf.Host, Path: "/login"}
	v := url.Values{}
	rvalue := rv{}

	if C.Token == "" {
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
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(body, &rvalue)
			if err != nil {
				panic(err)
			}
		}

		token := (rvalue.Data).(string)
		// 更新token, 写入token到文件
		err = C.UpdateToken(token)
		if err != nil {
			os.Exit(1)
		}
	} else {
		u.Path = "/api/v1/auth"
		client := &http.Client{}
		req, _ := http.NewRequest("GET", u.String(), nil)
		req.Header.Add("Token", C.Token)
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
		if err != nil {
			panic(err)
		}
		if rvalue.Code != 0 {
			log.ErrorLogger.Println("token invalid")
			C.Token = ""
			login()
		}
	}
}
