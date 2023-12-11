package constant

import (
	"os"
	"strings"

	"github.com/dmzlingyin/clipshare/pkg/log"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	UserName string `yaml:"username"`
	PassWord string `yaml:"password"`
	Device   string `yaml:"device"`
	Host     string `yaml:"host"`
	Mute     bool   `yaml:"mute"`
}

var Config Conf

func init() {
	data := readFile("./conf/client.yaml")
	err := yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Error.Fatalln(err)
	}
}

func readFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		if strings.Contains(path, "token") {
			return []byte{}
		}
		log.Error.Println(path, "open fail")
		panic(err)
	}
	return data
}
