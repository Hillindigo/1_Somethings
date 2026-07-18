package main

import (
	"blogx_server/service/text_service"
	"fmt"
	"os"
)

func main() {
	byteData, _ := os.ReadFile("text.md")
	list := text_service.MdContentTransformation(1,
		"xxx",
		string(byteData),
	)
	fmt.Println(list)
}
