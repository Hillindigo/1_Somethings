package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {
	var role int8
	fmt.Println("选择角色     1 超级管理员   2 普通用户   3 访客")
	_, err := fmt.Scanln(&role)
	if err != nil {
		logrus.Errorf("输入错误 %s", err)
		return
	}
	fmt.Println(role)
}
