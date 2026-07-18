// core/init_conf.go
package core

import (
	"blogx_server/conf"
	"blogx_server/flags"
	"blogx_server/global"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
)

func ReadConf() (c *conf.Config) {
	byteData, err := os.ReadFile(flags.FlagOptions.File)
	if err != nil {
		panic(err)
	}
	c = new(conf.Config)
	err = yaml.Unmarshal(byteData, c)
	if err != nil {
		panic(fmt.Sprintf("yaml配置文件格式错误 %s", err))
	}

	fmt.Printf("读取配置文件 %s 成功\n", flags.FlagOptions.File)
	return
}

func SetConf() {
	byteData, err := yaml.Marshal(global.Config)
	if err != nil {
		logrus.Errorf("conf读取失败 %s", err)
		return
	}

	err = os.WriteFile(flags.FlagOptions.File, byteData, 0666)
	if err != nil {
		logrus.Errorf("设置配置文件失败 %s", err)
		return
	}
}
