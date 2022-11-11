//go:build client
// +build client

package hub

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	C "github.com/dmzlingyin/clipshare/pkg/constant"
	"github.com/dmzlingyin/clipshare/pkg/log"
)

type CData struct {
	Data []byte
}

type rv struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Send(token string, data []byte) {
	cdata := CData{data}
	mdata, err := json.Marshal(cdata)
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}

	u := url.URL{Scheme: "http", Host: C.ClientConf.Host, Path: "/api/v1/transfer"}
	client := &http.Client{}
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(mdata))
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
	req.Header.Add("Token", token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	var rvalue rv
	err = json.Unmarshal(body, &rvalue)
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
	if rvalue.Code != 0 {
		log.ErrorLogger.Println(rvalue.Msg)
	}
}
