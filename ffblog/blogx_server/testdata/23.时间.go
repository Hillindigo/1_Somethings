package main

import (
	"fmt"
	"time"
)

func ywt() {
	now := time.Now()
	endTime := now.Format("2006-01-02") + " 23:59:59"
	endTimeObj, _ := time.Parse("2006-01-02 15:04:05", endTime)
	subTime := endTimeObj.Sub(now)
	fmt.Println(subTime.Seconds(), now, endTimeObj)
}

func zc() {
	// 获取当前时间
	now := time.Now()

	// 获取当前时区
	location := time.Local

	// 设置今天的结束时间为23:59:59，基于当前时区
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, location)

	// 打印今天的结束时间
	fmt.Println(endTime)
	subTime := endTime.Sub(now)
	fmt.Println(subTime.Seconds())
}

func main() {
	zc()
}
