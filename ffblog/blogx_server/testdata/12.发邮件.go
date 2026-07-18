// testdata/12.发邮件.go
package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/email_service"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()

	email_service.SendRegisterCode("2663456523@qq.com", "5433")
}
