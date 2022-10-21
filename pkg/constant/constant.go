package constant

import (
	"os"

	"github.com/dmzlingyin/clipshare/pkg/log"
	"gopkg.in/yaml.v2"
)

type CConf struct {
	IsRegister bool
	UserName   string
	PassWord   string
	Device     string
}

var ClientConf = CConf{}

func init() {
	data, err := os.ReadFile("./conf/client.yaml")
	if err != nil {
		log.ErrorLogger.Fatalln(err)
	}

	err = yaml.Unmarshal(data, &ClientConf)
	if err != nil {
		log.ErrorLogger.Fatalln(err)
	}
}
