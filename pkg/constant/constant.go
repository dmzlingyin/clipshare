package constant

import (
	"os"

	"github.com/dmzlingyin/clipshare/pkg/log"
	"gopkg.in/yaml.v2"
)

// 客户端配置
type CConf struct {
	IsRegister bool
	UserName   string
	PassWord   string
	Device     string
	Host       string
}

// 服务端配置
type SConf struct {
	MaxUsers   int
	MaxDevices int
}

var (
	ClientConf = CConf{}
	ServerConf = SConf{}
)

func init() {
	data, err := os.ReadFile("./conf/client.yaml")
	if err != nil {
		log.ErrorLogger.Fatalln("client.yaml open fail")
	}

	err = yaml.Unmarshal(data, &ClientConf)
	if err != nil {
		log.ErrorLogger.Fatalln(err)
	}

	data, err = os.ReadFile("./conf/server.yaml")
	if err != nil {
		log.ErrorLogger.Fatalln("server.yaml open fail")
	}

	err = yaml.Unmarshal(data, &ServerConf)
	if err != nil {
		log.ErrorLogger.Fatalln(err)
	}
}
