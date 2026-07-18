package main

import (
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

	//err := global.DB.Create(&models.ArticleModel{
	//	Title:   "嘻嘻嘻",
	//	TagList: ctype.List{"python", "go"},
	//}).Error
	//fmt.Println(err)
	var list1 []models.ArticleModel
	global.DB.Find(&list1)
	fmt.Println(list1)
}
