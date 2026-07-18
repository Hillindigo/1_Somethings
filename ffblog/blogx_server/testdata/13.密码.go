// testdata/13.密码.go
package main

import (
	"blogx_server/utils/pwd"
	"fmt"
)

func main() {
	p, _ := pwd.GenerateFromPassword("123456")
	fmt.Println(p)
	// $2a$10$HBN7dx4PPxzIsGyoS9RT.eAkyW.nwZPbz2gOYBHIq7A7BB06yqqXK

	fmt.Println(pwd.CompareHashAndPassword(p, "123456"))

}
