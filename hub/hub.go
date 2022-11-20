//go:build client
// +build client

package hub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	C "github.com/dmzlingyin/clipshare/pkg/constant"
)

func Send(token string, data []byte) error {
	cdata := CD{data}
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
	var rv RV
	err = json.Unmarshal(body, &rv)
	if err != nil {
		return err
	}
	if rv.Code != 0 {
		return errors.New("fail")
	}
	return nil
}

func SendFile(fp string) {
	u := url.URL{Scheme: "http", Host: C.ClientConf.Host, Path: "/api/v1/backup"}
	resp, err := http.PostForm(u.String(), nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
