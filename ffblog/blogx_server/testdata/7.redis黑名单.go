// testdata/7.redis黑名单.go
package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/redis_service/redis_jwt"
	"fmt"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Redis = core.InitRedis()

	//token, err := jwts.GetToken(jwts.Claims{
	//	UserID: 2,
	//	Role:   1,
	//})
	//fmt.Println(token, err)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjEsInVzZXJuYW1lIjoiIiwicm9sZSI6MSwiZXhwIjoxNzI5ODY3MzgwLCJpc3MiOiJmZW5nZmVuZyJ9.Nu2GwsfVFXYyh4EjU0Pos-D5plx5055ivIbY5JxBfng"
	redis_jwt.TokenBlack(token, redis_jwt.UserBlackType)
	blk, ok := redis_jwt.HasTokenBlack(token)
	fmt.Println(blk, ok)
}
