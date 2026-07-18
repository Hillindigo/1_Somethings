package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/chat_service"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.DB = core.InitDB()

	//fmt.Println(focus_service.CalcUserRelationship(1, 2))
	//fmt.Println(focus_service.CalcUserRelationship(1, 3))
	//fmt.Println(focus_service.CalcUserRelationship(4, 1))
	//fmt.Println(focus_service.CalcUserRelationship(5, 1))
	//fmt.Println(focus_service.CalcUserRelationship(3, 4))
	//fmt.Println(focus_service.CalcUserPatchRelationship(1, []uint{2, 3})) // 2：4  3：2   4：2  5:3
	//fmt.Println(focus_service.CalcUserPatchRelationship(2, []uint{1, 3, 4, 5}))
	//fmt.Println(focus_service.CalcUserPatchRelationship(3, []uint{1, 2, 4, 5}))
	//fmt.Println(focus_service.CalcUserPatchRelationship(4, []uint{1, 2, 3, 5}))
	//fmt.Println(focus_service.CalcUserPatchRelationship(5, []uint{1, 2, 3, 4}))

	chat_service.ToImageChat(2, 1, "http://image.fengfengzhidao.com/gvb_1009/20231122175551__1122c.png")
	chat_service.ToImageChat(1, 2, "http://image.fengfengzhidao.com/gvb_1009/20231122175551__1122c.png")
	chat_service.ToImageChat(1, 2, "http://image.fengfengzhidao.com/gvb_1009/20231122175551__1122c.png")

}
