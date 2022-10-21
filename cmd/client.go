package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

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
		// 将命令行传递的第一个参数作为用户名, 第二个参数作为device
		client(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	err := clipboard.Init()
	if err != nil {
		log.ErrorLogger.Fatalf("clipboard init failed")
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func client(username, device string) {
	fmt.Println("user: ", username, "connecting to server...")
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/socket"}

	// 建立websocket连接, 通过header区分客户端
	header := http.Header{}
	header.Add("UserName", username)
	header.Add("Device", device)
	c, _, err := (&websocket.Dialer{}).Dial(u.String(), header)
	if err != nil {
		log.ErrorLogger.Println(err)
	}
	defer c.Close()

	go Watch(username, device)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.ErrorLogger.Println(err)
			return
		}
		if len(message) > 0 {
			fmt.Println(username, device, string(message))
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
	uri := "http://localhost:8080/transfer"

	_, err := http.Post(uri, "application/json", bytes.NewReader(marshalForm))
	if err != nil {
		log.ErrorLogger.Println(err)
	}
}

func Watch(username, device string) {
	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	for data := range ch {
		fmt.Println(string(data))
		send(username, device, data)
	}
}
