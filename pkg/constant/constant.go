package constant

import (
	"fmt"
	"os"

	"github.com/dmzlingyin/clipshare/pkg/log"
	"gopkg.in/yaml.v2"
)

// 客户端配置
type CConf struct {
	UserName string `yaml:"username"`
	PassWord string `yaml:"password"`
	Device   string `yaml:"device"`
	Host     string `yaml:"host"`
}

// 服务端配置
type SConf struct {
	MaxUsers   int `yaml:"max_users"`
	MaxDevices int `yaml:"max_devices"`
	Port       int `yaml:"port"`
}

var (
	ClientConf = CConf{}
	ServerConf = SConf{}
	Token      = ""
)

func init() {
	data := readFile("./conf/client.yaml")
	err := yaml.Unmarshal(data, &ClientConf)
	if err != nil {
		log.ErrorLogger.Fatalln(err)
	}

	data = readFile("./conf/server.yaml")
	err = yaml.Unmarshal(data, &ServerConf)
	if err != nil {
		log.ErrorLogger.Fatalln(err)
	}

	data = readFile("./conf/.token")
	Token = string(data)
	fmt.Println(Token)
}

func readFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.ErrorLogger.Println(path, "open fail")
		panic(err)
	}
	return data
}

func UpdateToken(token string) error {
	err := os.WriteFile("./conf/.token", []byte(token), 0666)
	if err != nil {
		log.ErrorLogger.Println("write token to file fail")
		return err
	}
	return nil
}
