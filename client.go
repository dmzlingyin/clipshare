package main

import (
	"context"
	"fmt"
	C "github.com/dmzlingyin/clipshare/pkg/constant"
	"github.com/dmzlingyin/clipshare/pkg/log"
	"github.com/dmzlingyin/clipshare/pkg/utils"
	"golang.design/x/clipboard"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

type rv struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 用于控制剪贴板监控开启或关闭
var ctx context.Context
var cancel context.CancelFunc

func init() {
	err := clipboard.Init()
	if err != nil {
		log.Error.Fatalf("clipboard init failed")
	}

	// 开启剪贴板监控
	ctx, cancel = context.WithCancel(context.Background())
	go watch(ctx)
	// 注册热键
	go mainthread.Init(fn)
}

func fn() {
	// Ctrl + Shift + S
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
	err := hk.Register()
	if err != nil {
		log.Error.Println(err)
	}

	for {
		<-hk.Keydown()
		// 文件分享
		fmt.Println("hotkey down")
	}
}

func client() {
	for {
		//_, message, err := c.ReadMessage()
		var err error
		var message []byte
		if err != nil {
			log.Error.Println(err)
			return
		}
		if len(message) > 0 {
			// 关闭剪贴板监控
			cancel()
			clipboard.Write(clipboard.FmtText, message)
			if !C.ClientConf.Mute {
				go utils.Play("./docs/receive.wav")
			}
			// 重新开启剪贴板监控
			ctx, cancel = context.WithCancel(context.Background())
			go watch(ctx)
		}
	}
}

func watch(ctx context.Context) {
	ch := clipboard.Watch(ctx, clipboard.FmtText)
	for data := range ch {
		if len(data) != 0 {
			//err := hub.Send(token, data)
			var err error
			if err != nil {
				log.Error.Println(err)
			} else {
				if !C.ClientConf.Mute {
					// 播放发送成功提示音
					go utils.Play("./docs/send.wav")
				}
			}
		}
	}
}
