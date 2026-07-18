package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Join([]string{}, ","))
	fmt.Println(strings.Join([]string{"python"}, ","))
	fmt.Println(strings.Join([]string{"python", "go"}, ","))
}
