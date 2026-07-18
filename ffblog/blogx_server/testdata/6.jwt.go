// testdata/6.jwt.go
package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/utils/jwts"
	"fmt"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	token, err := jwts.GetToken(jwts.Claims{
		UserID: 1,
		Role:   1,
	})
	fmt.Println(token, err)
	//cls, err := jwts.ParseToken("xx")
	//fmt.Println(cls, err)
}
