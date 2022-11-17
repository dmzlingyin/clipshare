//go:build client
// +build client

package hub

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	C "github.com/dmzlingyin/clipshare/pkg/constant"
)

type CData struct {
	Data []byte
}

type rv struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Send(token string, data []byte) error {
	cdata := CData{data}
	mdata, err := json.Marshal(cdata)
	if err != nil {
		return nil
	}

	u := url.URL{Scheme: "http", Host: C.ClientConf.Host, Path: "/api/v1/transfer"}
	client := &http.Client{}
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(mdata))
	if err != nil {
		return err
	}
	req.Header.Add("Token", token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	var rvalue rv
	err = json.Unmarshal(body, &rvalue)
	if err != nil {
		return err
	}
	if rvalue.Code != 0 {
		return errors.New("fail")
	}
	return nil
}
