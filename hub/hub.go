//go:build client
// +build client

package hub

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	C "github.com/dmzlingyin/clipshare/pkg/constant"
	"github.com/dmzlingyin/clipshare/pkg/log"
)

type Meta struct {
	UserName string
	Device   string
	Data     []byte
}

func Send(username, device string, data []byte) {
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
