package main

import (
	"blogx_server/utils/computer"
	"fmt"
)

func main() {
	fmt.Println(computer.GetCpuPercent())
	fmt.Println(computer.GetMemPercent())
	fmt.Println(computer.GetDiskPercent())
}
