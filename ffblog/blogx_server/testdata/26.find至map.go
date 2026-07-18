package main

import (
	"blogx_server/common"
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models"
	"fmt"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.DB = core.InitDB()

	//mps := common.ScanMap(models.UserModel{}, common.ScanOption{
	//	Where: global.DB.Where("id in (1)"),
	//})
	mps := common.ScanMapV2(models.ChatModel{}, common.ScanOption{
		Where: global.DB.Where("id in (5,6,7)"),
	})
	fmt.Println(mps)
}
