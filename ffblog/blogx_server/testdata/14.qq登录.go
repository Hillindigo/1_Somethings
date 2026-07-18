// testdata/14.qq登录.go
package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/qq_service"
	"fmt"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	fmt.Println(qq_service.GetUserInfo("C7C4360BE86B5B91F5CFE3A0ECC6574C"))
}
