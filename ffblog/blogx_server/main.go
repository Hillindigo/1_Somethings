// /main.go
package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/router"
	"blogx_server/service/cron_service"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	core.InitIPDB()
	global.DB = core.InitDB()
	global.Redis = core.InitRedis()
	global.ESClient = core.EsConnect()

	flags.Run()

	core.InitMysqlES()

	// 定时任务
	cron_service.Cron()

	// 启动web程序
	router.Run()
}
