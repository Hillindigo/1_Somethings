// testdata/10.七牛云.go
package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/qiniu_service"
	"fmt"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	//url, err := SendFile("uploads/images/头像_0003_26.jpg")
	//fmt.Println(url, err)

	//file, _ := os.Open("uploads/images/头像_0003_26.jpg")
	//
	//url, err := SendReader(file)
	//fmt.Println(url, err)
	//fmt.Println(qiniu_service.GenToken())
	fmt.Println(qiniu_service.SendFile("./testdata/20231006164747__宇宙.png"))
}
