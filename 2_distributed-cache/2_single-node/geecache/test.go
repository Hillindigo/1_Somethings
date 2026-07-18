package main

import (
	"fmt"
	"sync"
	"time"
)

var set = make(map[int]bool, 0)

var l sync.Mutex

func printOnce(num int) {
	l.Lock()
	defer l.Unlock()
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true

}

func main() {
	for i := 0; i < 10; i++ {
		go printOnce(100)
	}
	time.Sleep(time.Second)
}
